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
	rawMap := make(map[types.RepoName]issueList)
	seen := make(map[int64]*github.Issue)
	for _, issue := range issues {
		if _, ok := seen[issue.GetID()]; ok {
			continue
		}
		seen[issue.GetID()] = issue
		url, err := url.Parse(issue.GetHTMLURL())
		if err != nil {
			return nil, err
		}
		path := strings.Split(url.Path, "/")
		if len(path) < 3 {
			return nil, fmt.Errorf("issue url path too short: %s", url.Path)
		}
		n := types.RepoName(path[2])
		rawMap[n] = append(rawMap[n], issue)
	}
	result := make(map[types.RepoName][]types.MyIssue)
	for n, v := range rawMap {
		lst := make([]types.MyIssue, len(v))
		for i, x := range sortIssuesByDateOfUpdate(v) {
			lst[i] = types.MyIssue{
				Number:  x.GetNumber(),
				Title:   x.GetTitle(),
				HtmlUrl: x.GetHTMLURL(),
				Updated: x.GetUpdatedAt().Time,
			}
		}
		result[n] = lst
	}
	return result, nil
}

func makeMapOfRepoToCommitList(commits commitList) (map[types.RepoName][]types.MyIssue, error) {
	rawMap := make(map[types.RepoName]commitList)
	seen := make(map[string]*github.CommitResult)
	for _, commit := range commits {
		if _, ok := seen[commit.GetSHA()]; ok {
			continue
		}
		seen[commit.GetSHA()] = commit
		url, err := url.Parse(commit.GetHTMLURL())
		if err != nil {
			return nil, err
		}
		path := strings.Split(url.Path, "/")
		if len(path) < 3 {
			return nil, fmt.Errorf("commit url path too short: %s", url.Path)
		}
		n := types.RepoName(path[2])
		rawMap[n] = append(rawMap[n], commit)
	}
	result := make(map[types.RepoName][]types.MyIssue)
	for n, v := range rawMap {
		lst := make([]types.MyIssue, len(v))
		for i, x := range sortCommitsByDateOfCommit(v) {
			lst[i] = types.MyIssue{
				Number:  0,
				Title:   upToFirstLfOrEnd(x.GetCommit().GetMessage()),
				HtmlUrl: x.GetHTMLURL(),
				Updated: x.GetCommit().GetCommitter().GetDate().Time,
			}
		}
		result[n] = lst
	}
	return result, nil
}

func upToFirstLfOrEnd(s string) string {
	if i := strings.Index(s, "\n"); i > 0 {
		return s[0:i]
	}
	return s
}

func removeDuplicateCommits(commits commitList) commitList {
	var noDupes commitList
	seen := make(map[string]*github.CommitResult)
	for _, commit := range commits {
		if _, ok := seen[commit.GetSHA()]; !ok {
			seen[commit.GetSHA()] = commit
			noDupes = append(noDupes, commit)
		}
	}
	return noDupes
}

func sortIssuesByDateOfUpdate(list issueList) issueList {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt().After(list[j].GetUpdatedAt().Time)
	})
	return list
}

func sortCommitsByDateOfCommit(list commitList) commitList {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetCommit().GetCommitter().GetDate().After(list[j].GetCommit().GetCommitter().GetDate().Time)
	})
	return list
}

func sortPrsDate(list []*github.PullRequest) []*github.PullRequest {
	sort.Slice(list, func(i, j int) bool {
		return list[i].MergedAt.After(list[j].MergedAt.Time)
	})
	return list
}
