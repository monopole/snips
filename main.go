package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"snips/internal"
	"sort"
	"strings"
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
	q.reportOrgs()
	q.reportUser()
	//fmt.Print("## Theme\n\n> _TODO_\n")
	//q.reportIssuesFiled()
	//q.reportReviews()
	//q.reportPrs()
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
	q.printIssues("Reviews", results.Issues)
}

func (q *questioner) printIssues(title string, issues []*github.Issue) {
	fmt.Printf("\n## %s\n\n", title)
	for repo, issueList := range convertToIssueMap(issues) {
		fmt.Printf("#### %s\n\n", repo)
		for _, issue := range issueList {
			q.printIssueLink(issue)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (q *questioner) reportIssuesFiled() {
	results, _, err := q.client.Search.Issues(
		q.ctx,
		fmt.Sprintf(
			"author:%s created:%s..%s",
			q.user,
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
	q.printIssues("Issues filed or commented on", results.Issues)
}

// Alternative query:
//  author:monopole is:pr merged:2019-11-25..2019-12-01
func (q *questioner) reportPrs() {
	fmt.Print("\n## PRS\n\n")
	for _, repo := range myRepos {
		prList, _, err := q.client.PullRequests.List(
			q.ctx,
			repo.org,
			repo.name,
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
		prList = q.filterPrs(prList)
		if len(prList) > 0 {
			//prList = sortPrsDate(prList)
			fmt.Printf("#### %s\n\n", repo.name)
			for _, pr := range prList {
				q.printPrLink(pr)
			}
			fmt.Println()
		}
	}
}

func sortPrsDate(list []*github.PullRequest) []*github.PullRequest {
	sort.Slice(list, func(i, j int) bool {
		return list[i].MergedAt.After(*list[j].MergedAt)
	})
	return list
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

func (q *questioner) printPrLink(pr *github.PullRequest) {
	fmt.Printf(" - %s [%s](%s)\n",
		pr.GetMergedAt().Format("2006-01-02"), pr.GetTitle(), pr.GetHTMLURL())
}

func convertToIssueMap(issues []*github.Issue) map[string][]*github.Issue {
	almost := make(map[string][]*github.Issue)
	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}
		issueUrl, err := url.Parse(issue.GetHTMLURL())
		if err != nil {
			log.Fatal(err)
		}
		path := strings.Split(issueUrl.Path, "/")
		repo := path[2]
		var list []*github.Issue
		if oldList, ok := almost[repo]; ok {
			list = append(oldList, issue)
		} else {
			list = []*github.Issue{issue}
		}
		almost[repo] = list
	}
	result := make(map[string][]*github.Issue)
	for repo, issueList := range almost {
		result[repo] = sortIssuesByDate(issueList)
	}
	return result
}

func sortIssuesByDate(list []*github.Issue) []*github.Issue {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt().After(list[j].GetUpdatedAt())
	})
	return list
}

func (q *questioner) printIssueLink(issue *github.Issue) {
	fmt.Printf(
		" - %s [%s](%s)\n",
		issue.GetUpdatedAt().Format("2006-01-02"), *issue.Title, *issue.HTMLURL)
}

func (q *questioner) reportOrgs() {
		orgs, _, err := q.client.Organizations.List(
			q.ctx, "", &github.ListOptions{
				Page:    0,
				PerPage: 100,
			})
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
