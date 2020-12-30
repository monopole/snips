package internal

import (
	"sort"

	"github.com/google/go-github/v33/github"
)

func sortIssuesByDateOfUpdate(list []*github.Issue) []*github.Issue {
	sort.Slice(list, func(i, j int) bool {
		return list[i].GetUpdatedAt().After(list[j].GetUpdatedAt())
	})
	return list
}

func sortPrsDate(list []*github.PullRequest) []*github.PullRequest {
	sort.Slice(list, func(i, j int) bool {
		return list[i].MergedAt.After(*list[j].MergedAt)
	})
	return list
}
