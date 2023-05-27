package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal"
	"log"
	"time"
)

type questioner struct {
	users     []string
	dateStart time.Time
	dateEnd   time.Time
	client    *github.Client
	ctx       context.Context
}

func (q questioner) doIt() {
	for _, user := range q.users {
		q.reportUser(user)
		q.reportOrgs()

		fmt.Printf("## Theme for %s\n\n> _TODO_\n",
			internal.DateRange(q.dateStart, q.dateEnd))

		issues := q.searchIssues(
			"created", fmt.Sprintf("author:%s", user))
		internal.PrintIssues("Issues created", internal.RemovePrsFrom(issues))

		issues = q.searchIssues(
			"closed", fmt.Sprintf("assignee:%s", user))
		internal.PrintIssues("Issues fixed/closed", internal.RemovePrsFrom(issues))

		issues = q.searchIssues(
			"merged", fmt.Sprintf("author:%s", user))
		internal.PrintIssues("PRs merged", internal.KeepOnlyPrsFrom(issues))

		issues = q.searchIssues(
			"updated", fmt.Sprintf("-author:%s commenter:%s", user, user))
		issues = append(issues,
			q.searchIssues("updated", fmt.Sprintf("reviewed-by:%s", user))...)
		internal.PrintIssues("Issues Commented", internal.RemovePrsFrom(issues))
		internal.PrintIssues("PRs Reviewed", internal.KeepOnlyPrsFrom(issues))
	}
}

func (q questioner) searchIssues(dateQualifier, baseQuery string) []*github.Issue {
	var issues []*github.Issue
	query := fmt.Sprintf(
		"%s:%s..%s %s",
		dateQualifier,
		q.dateStart.Format("2006-01-02"),
		q.dateEnd.Format("2006-01-02"),
		baseQuery)
	opts := internal.MakeSearchOptions()
	for {
		results, resp, err := q.client.Search.Issues(
			q.ctx,
			query,
			opts)
		if err != nil {
			log.Fatal(err)
		}
		issues = append(issues, results.Issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return issues
}

func (q questioner) reportOrgs() {
	opts := internal.MakeListOptions()
	orgs, _, err := q.client.Organizations.List(q.ctx, "", &opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("## Organizations\n")
	for i, organization := range orgs {
		fmt.Printf("  %v. %v\n", i+1, organization.GetLogin())
	}
}

func (q questioner) reportUser(u string) {
	user, _, err := q.client.Users.Get(q.ctx, u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("# %s %s\n", user.GetName(), user.GetCompany())
}
