package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"snips/internal"
	"time"

	"github.com/google/go-github/v33/github"
)

type questioner struct {
	user      string
	dateStart time.Time
	dateEnd   time.Time
	client    *github.Client
	ctx       context.Context
}

type ghRepo struct {
	org  string
	name string
}

// TODO: get this from command line
var myRepos = []ghRepo{
	{"kubernetes-sigs", "kustomize"},
	{"kubernetes-sigs", "cli-utils"},
	{"GoogleContainerTools", "kpt"},
}

func main() {
	if !(len(os.Args) == 4 || len(os.Args) == 5) {
		fmt.Print(`usage:
  snips {user} {githubAuthToken} {dateStart} [{dayCount}] 
e.g.
  go run . monopole deadbeef0000deadbeef 2020-04-06 
`)
		os.Exit(1)
	}
	user := os.Args[1]
	token := os.Args[2]
	dayStart := internal.ParseDate(os.Args[3])
	dayCount := 7
	if len(os.Args) == 5 {
		dayCount = internal.ParseDayCount(os.Args[4])
	}
	ctx := context.Background()
	q := questioner{
		user:      user,
		dateStart: dayStart,
		dateEnd:   dayStart.AddDate(0, 0, dayCount),
		ctx:       ctx,
		client:    internal.MakeClient(ctx, token),
	}
	//q.reportOrgs()
	//q.reportUser()
	fmt.Print("## Theme\n\n> _TODO_\n")

	//internal.PrintIssues("Issues filed or commented on", q.queryIssuesCreated())

	//internal.PrintIssues("Reviews", q.queryReviews())

	q.reportPrs()
}

func (q *questioner) queryReviews() (issues []*github.Issue) {
	opts := internal.MakeSearchOptions()
	for {
		results, resp, err := q.client.Search.Issues(
			q.ctx,
			// The query excludes the user as the author, looking
			// for the user only in comments.
			fmt.Sprintf(
				"-author:%s commenter:%s updated:%s..%s",
				q.user, q.user,
				q.dateStart.Format("2006-01-02"),
				q.dateEnd.Format("2006-01-02")),
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
	return
}

func (q *questioner) queryIssuesCreated() (issues []*github.Issue) {
	opts := internal.MakeSearchOptions()
	for {
		results, resp, err := q.client.Search.Issues(
			q.ctx,
			fmt.Sprintf(
				"author:%s created:%s..%s",
				q.user,
				q.dateStart.Format("2006-01-02"),
				q.dateEnd.Format("2006-01-02")),
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
	return
}

func (q *questioner) reportPrs() {
	fmt.Print("\n## PRS\n\n")
	for _, repo := range myRepos {
		prs := q.queryPrs(repo)
		if len(prs) > 0 {
			//prList = sortPrsDate(prList)
			fmt.Printf("#### %s\n\n", repo.name)
			for _, pr := range prs {
				internal.PrintPrLink(pr)
			}
			fmt.Println()
		}
	}
}

// Alternative query:
//  author:monopole is:pr merged:2019-11-25..2019-12-01
func (q *questioner) queryPrs(repo ghRepo) (prs []*github.PullRequest) {
	lOpts := internal.MakeListOptions()
	for {
		prList, resp, err := q.client.PullRequests.List(
			q.ctx,
			repo.org,
			repo.name,
			&github.PullRequestListOptions{
				State:       "closed",
				Head:        q.user + ":master",
				Base:        "",
				Sort:        "created",
				Direction:   "desc",
				ListOptions: lOpts,
			})
		if err != nil {
			log.Fatal(err)
		}
		prs = append(prs, prList...)
		if resp.NextPage == 0 {
			break
		}
		fmt.Print(".")
		lOpts.Page = resp.NextPage
	}
	fmt.Println()
	return
}

func (q *questioner) filterPrs(list []*github.PullRequest) []*github.PullRequest {
	var result []*github.PullRequest
	for _, pr := range list {
		if pr.GetUser().GetLogin() == q.user /* TODO add date check */ {
			result = append(result, pr)
		}
	}
	return result
}

func (q *questioner) reportOrgs() {
	opts := internal.MakeListOptions()
	orgs, _, err := q.client.Organizations.List(q.ctx, "", &opts)
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
	fmt.Printf("%s %s\n", user.GetName(), user.GetCompany())
}
