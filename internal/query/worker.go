package query

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

// Worker queries GitHub.
type Worker struct {
	Users     []string
	DateRange *types.DayRange
	GhClient  *github.Client
	Ctx       context.Context
}

type issueList []*github.Issue
type commitList []*github.CommitResult

type filter int

const (
	unknown filter = iota
	keepOnlyPrs
	rejectPrs
)

func (w *Worker) DoIt() ([]*types.MyUser, error) {
	var result []*types.MyUser
	for _, n := range w.Users {
		time.Sleep(15 * time.Second) // Avoid hitting API rate limit.
		rec, err := w.doQueriesOnUser(n)
		if err != nil {
			log.Printf("trouble looking up user %s: %s\n", n, err.Error())
		}
		result = append(result, rec)
	}
	return result, nil
}

func (w *Worker) doQueriesOnUser(userName string) (*types.MyUser, error) {
	myUser, err := w.getUserRec(userName)
	if err != nil {
		return nil, err
	}
	if err = w.getOrgs(myUser); err != nil {
		return nil, err
	}
	if err = w.getIssues(myUser); err != nil {
		return nil, err
	}
	if err = w.getCommits(myUser); err != nil {
		return nil, err
	}
	return myUser, nil
}

func (w *Worker) getIssues(myUser *types.MyUser) (err error) {
	var lst issueList
	lst, err = w.searchIssues("created", "author:%s", myUser.Login)
	if err != nil {
		return err
	}
	myUser.IssuesCreated, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return err
	}

	lst, err = w.searchIssues("closed", "assignee:%s", myUser.Login)
	if err != nil {
		return err
	}
	myUser.IssuesClosed, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return err
	}

	lst, err = w.searchIssues("merged", "author:%s", myUser.Login)
	if err != nil {
		return err
	}
	myUser.PrsMerged, err = makeMapOfRepoToIssueList(filterIssues(lst, keepOnlyPrs))
	if err != nil {
		return err
	}

	lst, err = w.searchIssues("updated", "-author:%s commenter:%s", myUser.Login, myUser.Login)
	if err != nil {
		return err
	}
	{
		var lst2 issueList
		lst2, err = w.searchIssues("updated", "reviewed-by:%s", myUser.Login)
		if err != nil {
			return err
		}
		lst = append(lst, lst2...)
	}
	myUser.IssuesCommented, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	myUser.PrsReviewed, err = makeMapOfRepoToIssueList(filterIssues(lst, keepOnlyPrs))
	return nil
}

func (w *Worker) getUserRec(n string) (*types.MyUser, error) {
	var r types.MyUser
	user, _, err := w.GhClient.Users.Get(w.Ctx, n)
	if err != nil {
		log.Fatal(err)
	}
	r.Name = user.GetName()
	r.Company = user.GetCompany()
	r.Login = user.GetLogin()
	r.Email = user.GetEmail()
	return &r, nil
}

// searchIssues uses the "search" endpoint, not the "issues" endpoint, because the goal is to
// discover what the user has been doing with issues, rather than manage issues.
// Compare the docs here:
// https://docs.github.com/en/rest/search?apiVersion=2022-11-28#search-issues-and-pull-requests
// https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28
func (w *Worker) searchIssues(dateQualifier, qFmt string, args ...any) (issueList, error) {
	query := w.makeQuery(dateQualifier, qFmt, args)
	opts := makeSearchOptions()
	var lst issueList
	for {
		results, resp, err := w.GhClient.Search.Issues(w.Ctx, query, opts)
		if err != nil {
			return nil, err
		}
		lst = append(lst, results.Issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return lst, nil
}

// getCommits looks for non-merge commits
// (merge commits usually aren't interesting).
func (w *Worker) getCommits(myUser *types.MyUser) (err error) {
	var lst commitList
	lst, err = w.searchCommits("author-date", "merge:false author:%s", myUser.Login)
	if err != nil {
		return err
	}
	if lookAtCommitterField := false; lookAtCommitterField {
		// Usually the author is the committer, so don't bother?
		var lst2 commitList
		lst2, err = w.searchCommits("committer-date", "merge:false committer:%s", myUser.Login)
		if err != nil {
			return err
		}
		lst = append(lst, lst2...)
		lst = removeDuplicateCommits(lst)
	}
	myUser.Commits, err = makeMapOfRepoToCommitList(lst)
	return nil
}

// searchCommits uses https://docs.github.com/en/search-github/searching-on-github/searching-commits
// It doesn't use https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28
func (w *Worker) searchCommits(dateQualifier, qFmt string, args ...any) (commitList, error) {
	query := w.makeQuery(dateQualifier, qFmt, args)
	opts := makeSearchOptions()
	var lst commitList
	for {
		results, resp, err := w.GhClient.Search.Commits(w.Ctx, query, opts)
		if err != nil {
			return nil, err
		}
		lst = append(lst, results.Commits...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return lst, nil
}

func (w *Worker) makeQuery(dateQualifier string, qFmt string, args []any) string {
	return fmt.Sprintf(
		"%s:%s..%s %s",
		dateQualifier,
		w.DateRange.StartAsTime().Format(types.DayFormat1),
		w.DateRange.EndAsTime().Format(types.DayFormat1),
		fmt.Sprintf(qFmt, args...))
}

func filterIssues(issues issueList, f filter) (result issueList) {
	if f != keepOnlyPrs && f != rejectPrs {
		log.Fatalf("unable to deal with filter %v", f)
	}
	for _, issue := range issues {
		if issue.IsPullRequest() == (f == keepOnlyPrs) {
			result = append(result, issue)
		}
	}
	return
}

func (w *Worker) getOrgs(u *types.MyUser) error {
	lOpts := makeListOptions()
	orgs, _, err := w.GhClient.Organizations.List(w.Ctx, u.Login, &lOpts)
	if err != nil {
		return err
	}
	for i := range orgs {
		u.Orgs = append(u.Orgs, types.MyOrg{
			Name:  orgs[i].GetName(),
			Login: orgs[i].GetLogin(),
		})
	}
	return nil
}

func makeSearchOptions() *github.SearchOptions {
	return &github.SearchOptions{
		Sort:        "",
		Order:       "", // "asc", "desc"
		TextMatch:   false,
		ListOptions: makeListOptions(),
	}
}

func makeListOptions() github.ListOptions {
	return github.ListOptions{PerPage: 50}
}
