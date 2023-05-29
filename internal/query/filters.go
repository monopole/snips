package query

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

func makeMapOfRepoToIssueList(issues issueList) (map[types.RepoName][]types.MyIssue, error) {
	rawIssues := make(map[types.RepoName]issueList)
	seen := make(map[int64]*github.Issue)
	for _, issue := range issues {
		if _, ok := seen[issue.GetID()]; ok {
			continue
		}
		seen[*issue.ID] = issue
		issueUrl, err := url.Parse(issue.GetHTMLURL())
		if err != nil {
			return nil, err
		}
		path := strings.Split(issueUrl.Path, "/")
		if len(path) < 3 {
			return nil, fmt.Errorf("issueUrl path too short: %s", issueUrl.Path)
		}
		n := types.RepoName(path[2])
		rawIssues[n] = append(rawIssues[n], issue)
	}
	result := make(map[types.RepoName][]types.MyIssue)
	for n, v := range rawIssues {
		lst := make([]types.MyIssue, len(v))
		for i, x := range sortIssuesByDateOfUpdate(v) {
			lst[i] = types.MyIssue{
				Title:   x.GetTitle(),
				HtmlUrl: x.GetHTMLURL(),
				Updated: x.GetUpdatedAt().Time,
			}
		}
		result[n] = lst
	}
	return result, nil
}

func sortIssuesByDateOfUpdate(list issueList) issueList {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt().After(list[j].GetUpdatedAt().Time)
	})
	return list
}

func sortPrsDate(list []*github.PullRequest) []*github.PullRequest {
	sort.Slice(list, func(i, j int) bool {
		return list[i].MergedAt.After(list[j].MergedAt.Time)
	})
	return list
}
