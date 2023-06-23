package search

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

// Engine performs GitHub searches over a date range.
type Engine struct {
	ctx      context.Context
	client   *github.Client
	dayRange *types.DayRange
}

// MakeEngine returns an instance of a GitHub search engine.
func MakeEngine(ctx context.Context, cl *github.Client) *Engine {
	return &Engine{ctx: ctx, client: cl}
}

const pauseApiUs = 15 * time.Second

// LookupPeeps gathers data about the given usernames in the given day range.
func (se *Engine) LookupPeeps(names []string, dayRange *types.DayRange) ([]*types.MyUser, error) {
	var result []*types.MyUser
	se.dayRange = dayRange
	for i, n := range names {
		if i > 0 {
			// Avoid hitting API rate limit.
			time.Sleep(pauseApiUs)
		}
		fmt.Fprintf(os.Stderr, "Working on user %s...\n", n)
		if rec, err := se.doQueriesOnUser(n); err == nil {
			result = append(result, rec)
		} else {
			log.Printf("trouble with user %s: %s\n", n, err.Error())
		}
	}
	return result, nil
}

// searchIssues uses the "search" endpoint, not the "issues" endpoint, because the goal is to
// discover what the user has been doing with issues, rather than manage issues.
// https://docs.github.com/en/rest/search?apiVersion=2022-11-28#search-issues-and-pull-requests
// https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28
func (se *Engine) searchIssues(dateQualifier, qFmt string, args ...any) ([]*github.Issue, error) {
	query := se.makeQuery(dateQualifier, qFmt, args)
	opts := makeSearchOptions()
	var lst []*github.Issue
	for {
		results, resp, err := se.client.Search.Issues(se.ctx, query, opts)
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

// searchCommits uses https://docs.github.com/en/search-github/searching-on-github/searching-commits
// It doesn't use https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28
func (se *Engine) searchCommits(dateQualifier, qFmt string, args ...any) ([]*github.CommitResult, error) {
	query := se.makeQuery(dateQualifier, qFmt, args)
	opts := makeSearchOptions()
	var lst []*github.CommitResult
	for {
		results, resp, err := se.client.Search.Commits(se.ctx, query, opts)
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

func (se *Engine) makeQuery(dateQualifier string, qFmt string, args []any) string {
	return fmt.Sprintf(
		"%s:%s..%s %s",
		dateQualifier,
		se.dayRange.StartAsTime().Format(types.DayFormatGitHub),
		se.dayRange.EndAsTime().Format(types.DayFormatGitHub),
		fmt.Sprintf(qFmt, args...))
}

func (se *Engine) doQueriesOnUser(userName string) (*types.MyUser, error) {
	myUser, err := se.loadUserData(userName)
	if err != nil {
		return nil, err
	}
	if myUser.Orgs, err = se.findOrganizations(myUser); err != nil {
		return nil, err
	}
	lst, err := se.searchIssues("created", "author:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.IssuesCreated, err = makeMapOfRepoToIssueList(rejectPrs.from(lst))
	if err != nil {
		return nil, err
	}
	lst, err = se.searchIssues("closed", "assignee:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.IssuesClosed, err = makeMapOfRepoToIssueList(rejectPrs.from(lst))
	if err != nil {
		return nil, err
	}
	if myUser.IssuesCommented, myUser.PrsReviewed, err = se.findReviewsAndComments(myUser); err != nil {
		return nil, err
	}
	if myUser.Commits, err = se.findCommits(myUser); err != nil {
		return nil, err
	}
	return myUser, nil
}

func (se *Engine) loadUserData(n string) (*types.MyUser, error) {
	user, _, err := se.client.Users.Get(se.ctx, n)
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

func (se *Engine) findOrganizations(u *types.MyUser) ([]types.MyOrg, error) {
	lOpts := makeListOptions()
	orgs, _, err := se.client.Organizations.List(se.ctx, u.Login, &lOpts)
	if err != nil {
		return nil, err
	}
	var result []types.MyOrg
	for i := range orgs {
		result = append(result, types.MyOrg{
			Name:  orgs[i].GetName(),
			Login: orgs[i].GetLogin(),
		})
	}
	return result, nil
}

func (se *Engine) findReviewsAndComments(myUser *types.MyUser) (
	issuesReviewed map[types.RepoId][]types.MyIssue, prsReviewed map[types.RepoId][]types.MyIssue, err error) {
	lst, err := se.searchIssues("updated", "-author:%s commenter:%s", myUser.Login, myUser.Login)
	if err != nil {
		return
	}
	{
		var lst2 []*github.Issue
		lst2, err = se.searchIssues("updated", "reviewed-by:%s", myUser.Login)
		if err != nil {
			return
		}
		lst = append(lst, lst2...)
	}
	issuesReviewed, err = makeMapOfRepoToIssueList(rejectPrs.from(lst))
	if err != nil {
		return
	}
	prsReviewed, err = makeMapOfRepoToIssueList(keepOnlyPrs.from(lst))
	return
}
