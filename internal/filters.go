package internal

import (
	"github.com/google/go-github/v33/github"
	"log"
	"net/url"
	"strings"
)

func MapRepoToIssueList(issues []*github.Issue) map[string][]*github.Issue {
	almost := make(map[string][]*github.Issue)
	seen := make(map[int64]*github.Issue)
	for _, issue := range issues {
		if _, ok := seen[*issue.ID]; ok {
			continue
		} else {
			seen[*issue.ID] = issue
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
		result[repo] = SortIssuesByDateOfUpdate(issueList)
	}
	return result
}

func RemovePrsFrom(issues []*github.Issue) []*github.Issue {
	return filterIssue(issues, false)
}

func KeepOnlyPrsFrom(issues []*github.Issue) []*github.Issue {
	return filterIssue(issues, true)
}

func filterIssue(issues []*github.Issue, wantPr bool) (result []*github.Issue) {
	for _, issue := range issues {
		if issue.IsPullRequest() == wantPr {
			result = append(result, issue)
		}
	}
	return result
}
