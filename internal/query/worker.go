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
	if err = w.getAllCommitsInDateRange(myUser); err != nil {
		return nil, err
	}
	return myUser, nil
}

func (w *Worker) getIssues(myUser *types.MyUser) (err error) {
	var lst issueList
	if getCommentingActivity := false; getCommentingActivity {
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
	}
	if getIssuesToo := false; getIssuesToo {
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
	}
	if getPrsToo := false; getPrsToo {
		lst, err = w.searchIssues("merged", "author:%s", myUser.Login)
		if err != nil {
			return err
		}
		var prsMerged map[types.RepoId][]types.MyIssue
		if prsMerged, err = makeMapOfRepoToIssueList(filterIssues(lst, keepOnlyPrs)); err != nil {
			return err
		}
		commitMap := make(map[types.RepoId][]*types.MyCommit)
		for repoId, prList := range prsMerged {
			var allCommits []*types.MyCommit
			for i := range prList {
				var commitsForPr []*types.MyCommit
				commitsForPr, err = w.getCommitsForPr(&prList[i])
				if err != nil {
					return err
				}
				allCommits = append(allCommits, commitsForPr...)
			}
			commitMap[repoId] = allCommits
		}
		myUser.CommitsFromPrs = commitMap
	}
	return nil
}

func (w *Worker) getUserRec(n string) (*types.MyUser, error) {
	user, _, err := w.GhClient.Users.Get(w.Ctx, n)
	if err != nil {
		return nil, err
	}
	return &types.MyUser{
		Name:    user.GetName(),
		Company: user.GetCompany(),
		Login:   user.GetLogin(),
		Email:   user.GetEmail(),
	}, nil
}

// getCommitsForPr
// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#list-commits-on-a-pull-request
func (w *Worker) getCommitsForPr(prIssue *types.MyIssue) (result []*types.MyCommit, err error) {
	var (
		resp         *github.Response
		commits, lst []*github.RepositoryCommit
	)
	opts := makeListOptions()
	for {
		lst, resp, err = w.GhClient.PullRequests.ListCommits(
			w.Ctx, prIssue.RepoId.Org, prIssue.RepoId.Repo, prIssue.Number, &opts)
		if err != nil {
			return nil, err
		}
		commits = append(commits, lst...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	result = make([]*types.MyCommit, len(lst))
	for i, c := range lst {
		result[i] = &types.MyCommit{
			RepoId:           prIssue.RepoId,
			Sha:              c.GetSHA(),
			Url:              c.GetHTMLURL(),
			MessageFirstLine: upToFirstLfOrEnd(c.GetCommit().GetMessage()),
			Committed:        c.GetCommit().GetCommitter().GetDate().Time,
			Author:           c.GetAuthor().GetLogin(),
			Pr:               prIssue,
		}
	}
	return result, nil
}

// searchIssues uses the "search" endpoint, not the "issues" endpoint, because the goal is to
// discover what the user has been doing with issues, rather than manage issues.
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

// getAllCommitsInDateRange looks for non-merge commits
// (merge commits usually aren't interesting).
func (w *Worker) getAllCommitsInDateRange(myUser *types.MyUser) (err error) {
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
	myUser.CommitsBare, err = makeMapOfRepoToCommitList(lst)
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
