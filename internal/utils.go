package internal

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func ParseDate(arg string) time.Time {
	t, err := time.Parse("2006-01-02", arg)
	if err != nil {
		fmt.Printf("Trouble with date specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return t
}

func ParseDayCount(arg string) int {
	i, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Printf("Trouble with day count specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return i
}

func MakeClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func MapRepoToIssueList(issues []*github.Issue) map[string][]*github.Issue {
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
		result[repo] = sortIssuesByDateOfUpdate(issueList)
	}
	return result
}
