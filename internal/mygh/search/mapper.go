package search

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

type HasUrl interface {
	GetHTMLURL() string
}

func yankRepoId(raw HasUrl) (types.RepoId, error) {
	theUrl, err := url.Parse(raw.GetHTMLURL())
	if err != nil {
		return types.RepoId{}, err
	}
	path := strings.Split(theUrl.Path, "/")
	if len(path) < 3 {
		return types.RepoId{}, fmt.Errorf("issue url path too short in url: %q", raw)
	}
	return types.RepoId{
		Org:  path[1],
		Name: path[2],
	}, nil
}

func makeMapOfRepoToIssueList(issues []*github.Issue) (map[types.RepoId][]types.MyIssue, error) {
	rawMap := make(map[types.RepoId][]*github.Issue)
	seen := make(map[int64]*github.Issue)
	for _, issue := range issues {
		if _, ok := seen[issue.GetID()]; ok {
			continue
		}
		seen[issue.GetID()] = issue
		id, err := yankRepoId(issue)
		if err != nil {
			return nil, err
		}
		rawMap[id] = append(rawMap[id], issue)
	}
	result := make(map[types.RepoId][]types.MyIssue)
	for id, ghIssues := range rawMap {
		lst := make([]types.MyIssue, len(ghIssues))
		ghIssues = sortIssuesByDateOfUpdate(ghIssues)
		for i := range ghIssues {
			x := ghIssues[i]
			lst[i] = types.MyIssue{
				RepoId:  id,
				Number:  x.GetNumber(),
				Title:   x.GetTitle(),
				HtmlUrl: x.GetHTMLURL(),
				Updated: x.GetUpdatedAt().Time,
			}
		}
		result[id] = lst
	}
	return result, nil
}

func makeMapOfRepoToCommitList(commits []*github.CommitResult) (map[types.RepoId][]*types.MyCommit, error) {
	rawMap := make(map[types.RepoId][]*github.CommitResult)
	seen := make(map[string]*github.CommitResult)
	for _, commit := range commits {
		if _, ok := seen[commit.GetSHA()]; ok {
			continue
		}
		seen[commit.GetSHA()] = commit
		n, err := yankRepoId(commit)
		if err != nil {
			return nil, err
		}
		rawMap[n] = append(rawMap[n], commit)
	}
	result := make(map[types.RepoId][]*types.MyCommit)
	for id, ghCommits := range rawMap {
		lst := make([]*types.MyCommit, len(ghCommits))
		for i, x := range sortCommitsByDateOfCommit(ghCommits) {
			lst[i] = &types.MyCommit{
				RepoId:           id,
				Sha:              x.GetCommit().GetSHA(),
				Url:              x.GetHTMLURL(),
				MessageFirstLine: upToFirstLfOrEnd(x.GetCommit().GetMessage()),
				Committed:        x.GetCommit().GetCommitter().GetDate().Time,
				Author:           x.GetAuthor().GetLogin(),
				Pr:               nil, // This commit might not be associated with a Pr.
			}
		}
		result[id] = lst
	}
	return result, nil
}

func upToFirstLfOrEnd(s string) string {
	if i := strings.Index(s, "\n"); i > 0 {
		return s[0:i]
	}
	return s
}

func removeDuplicateCommits(commits []*github.CommitResult) []*github.CommitResult {
	var noDupes []*github.CommitResult
	seen := make(map[string]*github.CommitResult)
	for _, commit := range commits {
		if _, ok := seen[commit.GetSHA()]; !ok {
			seen[commit.GetSHA()] = commit
			noDupes = append(noDupes, commit)
		}
	}
	return noDupes
}

func sortIssuesByDateOfUpdate(list []*github.Issue) []*github.Issue {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt().After(list[j].GetUpdatedAt().Time)
	})
	return list
}

func sortCommitsByDateOfCommit(list []*github.CommitResult) []*github.CommitResult {
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
