package query

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

type issueList []*github.Issue

type Worker struct {
	Users     []string
	DateRange *types.DayRange
	GhClient  *github.Client
	Ctx       context.Context
}

type filter int

const (
	unknown filter = iota
	keepOnlyPrs
	rejectPrs
)

func (q Worker) DoIt() ([]*types.MyUser, error) {
	var result []*types.MyUser
	for _, n := range q.Users {
		rec, err := q.doQueriesOnUser(n)
		if err != nil {
			log.Printf("trouble looking up user %s: %s\n", n, err.Error())
		}
		result = append(result, rec)
	}
	return result, nil
}

func (q Worker) doQueriesOnUser(userName string) (*types.MyUser, error) {
	myUser, err := q.getUserRec(userName)
	if err != nil {
		return nil, err
	}
	err = q.getOrgs(myUser)
	if err != nil {
		return nil, err
	}

	var lst issueList

	lst, err = q.searchIssues("created", "author:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.IssuesCreated, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return nil, err
	}

	lst, err = q.searchIssues("closed", "assignee:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.IssuesClosed, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	if err != nil {
		return nil, err
	}
	lst, err = q.searchIssues("merged", "author:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	myUser.PrsMerged, err = makeMapOfRepoToIssueList(filterIssues(lst, keepOnlyPrs))
	if err != nil {
		return nil, err
	}
	lst, err = q.searchIssues("updated", "-author:%s commenter:%s", myUser.Login, myUser.Login)
	if err != nil {
		return nil, err
	}
	{
		var lst2 issueList
		lst2, err = q.searchIssues("updated", "reviewed-by:%s", myUser.Login)
		if err != nil {
			return nil, err
		}
		lst = append(lst, lst2...)
	}
	myUser.IssuesCommented, err = makeMapOfRepoToIssueList(filterIssues(lst, rejectPrs))
	myUser.PrsReviewed, err = makeMapOfRepoToIssueList(filterIssues(lst, keepOnlyPrs))
	return myUser, nil
}

func (q Worker) getUserRec(n string) (*types.MyUser, error) {
	var r types.MyUser
	user, _, err := q.GhClient.Users.Get(q.Ctx, n)
	if err != nil {
		log.Fatal(err)
	}
	r.Name = user.GetName()
	r.Company = user.GetCompany()
	r.Login = user.GetLogin()
	r.Email = user.GetEmail()
	return &r, nil
}

func (q Worker) searchIssues(dateQualifier, qFmt string, args ...any) (issueList, error) {
	query := fmt.Sprintf(
		"%s:%s..%s %s",
		dateQualifier,
		q.DateRange.StartAsTime().Format(types.DayFormat1),
		q.DateRange.EndAsTime().Format(types.DayFormat1),
		fmt.Sprintf(qFmt, args...))
	var lst issueList
	opts := makeSearchOptions()
	for {
		results, resp, err := q.GhClient.Search.Issues(q.Ctx, query, opts)
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

func (q Worker) getOrgs(r *types.MyUser) error {
	lOpts := makeListOptions()
	orgs, _, err := q.GhClient.Organizations.List(q.Ctx, r.Login, &lOpts)
	if err != nil {
		return err
	}
	for i := range orgs {
		r.Orgs = append(r.Orgs, types.MyOrg{
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
