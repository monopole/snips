package query

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/search"
	"github.com/monopole/snips/internal/types"
)

// Worker queries GitHub.
type Worker struct {
	Users    []string
	Se       *search.Engine
	GhClient *github.Client
	Ctx      context.Context
}

type filter int

const (
	unknown filter = iota
	keepOnlyPrs
	rejectPrs
	pauseApiUs = 15 * time.Second
)

func (w *Worker) DoIt() ([]*types.MyUser, error) {
	var result []*types.MyUser
	for i, n := range w.Users {
		if i > 0 {
			// Avoid hitting API rate limit.
			time.Sleep(pauseApiUs)
		}
		rec, err := w.doQueriesOnUser(n)
		if err != nil {
			log.Printf("trouble looking up user %s: %s\n", n, err.Error())
		}
		result = append(result, rec)
	}
	return result, nil
}

func (w *Worker) doQueriesOnUser(userName string) (*types.MyUser, error) {
	myUser, err := w.loadUserData(userName)
	if err != nil {
		return nil, err
	}
	if myUser.Orgs, err = w.findOrganizations(myUser); err != nil {
		return nil, err
	}
	lst, err := w.Se.SearchIssues("created", "author:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.IssuesCreated, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return nil, err
	}

	lst, err = w.Se.SearchIssues("closed", "assignee:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.IssuesClosed, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return nil, err
	}
	if myUser.IssuesCommented, myUser.PrsReviewed, err = w.findReviewsAndComments(myUser); err != nil {
		return nil, err
	}
	if myUser.Commits, err = w.findTheCommits(myUser); err != nil {
		return nil, err
	}
	return myUser, nil
}

func (w *Worker) findReviewsAndComments(myUser *types.MyUser) (
	issuesReviewed map[types.RepoId][]types.MyIssue, prsReviewed map[types.RepoId][]types.MyIssue, err error) {
	lst, err := w.Se.SearchIssues("updated", "-author:%s commenter:%s", myUser.Login, myUser.Login)
	if err != nil {
		return
	}
	{
		var lst2 []*github.Issue
		lst2, err = w.Se.SearchIssues("updated", "reviewed-by:%s", myUser.Login)
		if err != nil {
			return
		}
		lst = append(lst, lst2...)
	}
	issuesReviewed, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return
	}
	prsReviewed, err = makeMapOfRepoToIssueList(filterIssues(lst, keepOnlyPrs))
	return
}

func filterIssues(issues []*github.Issue, f filter) (result []*github.Issue) {
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

func (w *Worker) findOrganizations(u *types.MyUser) ([]types.MyOrg, error) {
	lOpts := search.MakeListOptions()
	orgs, _, err := w.GhClient.Organizations.List(w.Ctx, u.Login, &lOpts)
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

func (w *Worker) loadUserData(n string) (*types.MyUser, error) {
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
