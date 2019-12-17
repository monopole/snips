package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

type questioner struct {
	ghOrg     string
	ghRepo    string
	user      string
	dateStart time.Time
	dateEnd   time.Time
	client    *github.Client
	ctx       context.Context
}

func main() {
	if len(os.Args) < 7 {
		fmt.Print(`usage:
  snips {githubOrg} {githubRepo} {user} {dateStart} {dateEnd} {githubAuthToken}
e.g.
  snips kubernetes-sigs kustomize monopole 2019-11-25 2019-12-01 deadbeef0000deadbeef
`)
		os.Exit(1)
	}
	ctx := context.Background()
	q := questioner{
		ghOrg:     os.Args[1],
		ghRepo:    os.Args[2],
		user:      os.Args[3],
		dateStart: parseDate(os.Args[4]),
		dateEnd:   parseDate(os.Args[5]),
		ctx:       ctx,
		client:    makeClient(ctx, os.Args[6]),
	}
	//q.reportOrgs(client)
	//q.reportUser(client, ctx)
	q.reportReviews()
	q.reportPrs()
}

func parseDate(arg string) time.Time {
	t, err := time.Parse("2006-01-02", arg)
	if err != nil {
		fmt.Printf("Trouble with date specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return t
}

func makeClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

// The query excludes the user as the author, looking
// for the user only in comments.
func (q *questioner) reportReviews() {
	results, _, err := q.client.Search.Issues(
		q.ctx,
		fmt.Sprintf(
			"-author:%s commenter:%s updated:%s..%s",
			q.user, q.user,
			q.dateStart.Format("2006-01-02"),
			q.dateEnd.Format("2006-01-02")),
		&github.SearchOptions{
			Sort:      "",
			Order:     "",
			TextMatch: false,
			ListOptions: github.ListOptions{
				Page:    0,
				PerPage: 100,
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("## Reviews")
	for _, r := range results.Issues {
		q.printPrLink(r.GetTitle(), r.GetNumber(), r.GetUpdatedAt())
	}
}

// Alternative query:
//  author:monopole is:pr merged:2019-11-25..2019-12-01
func (q *questioner) reportPrs() {
	results, _, err := q.client.PullRequests.List(
		q.ctx,
		q.ghOrg,
		q.ghRepo,
		&github.PullRequestListOptions{
			State:     "closed",
			Head:      "",
			Base:      "",
			Sort:      "",
			Direction: "desc",
			ListOptions: github.ListOptions{
				Page:    0,
				PerPage: 100,
			},
		})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("## PRS")
	for _, r := range results {
		if r.GetUser().GetLogin() == q.user /* TODO add date check */ {
			q.printPrLink(
				r.GetTitle(), r.GetNumber(), r.GetMergedAt())
		}
	}
}

func (q *questioner) printPrLink(title string, prNum int, date time.Time) {
	fmt.Printf(
		" - %s [%s](https://github.com/%s/%s/pull/%d)\n",
		date.Format("2006-01-02"), title, q.ghOrg, q.ghRepo, prNum)
}

func (q *questioner) reportOrgs() {
	orgs, _, err := q.client.Organizations.List(
		q.ctx, q.ghOrg, nil)
	if err != nil {
		log.Fatal(err)
	}
	for i, organization := range orgs {
		fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	}
}

func (q *questioner) reportUser() {
	user, _, err := q.client.Users.Get(q.ctx, q.user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", user.GetName())
}
