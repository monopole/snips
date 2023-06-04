package search

import (
	"context"
	"fmt"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

// Engine performs GitHub searches over a date range.
type Engine struct {
	ctx       context.Context
	ghClient  *github.Client
	dateRange *types.DayRange
}

func MakeEngine(ctx context.Context, ghClient *github.Client, dateRange *types.DayRange) *Engine {
	return &Engine{
		ctx:       ctx,
		ghClient:  ghClient,
		dateRange: dateRange,
	}
}

// SearchIssues uses the "search" endpoint, not the "issues" endpoint, because the goal is to
// discover what the user has been doing with issues, rather than manage issues.
// https://docs.github.com/en/rest/search?apiVersion=2022-11-28#search-issues-and-pull-requests
// https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28
func (se *Engine) SearchIssues(dateQualifier, qFmt string, args ...any) ([]*github.Issue, error) {
	query := se.makeQuery(dateQualifier, qFmt, args)
	opts := makeSearchOptions()
	var lst []*github.Issue
	for {
		results, resp, err := se.ghClient.Search.Issues(se.ctx, query, opts)
		if err != nil {
			return nil, err
		}
		lst = append(lst, results.Issues...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return lst, nil
}

// SearchCommits uses https://docs.github.com/en/search-github/searching-on-github/searching-commits
// It doesn't use https://docs.github.com/en/rest/commits/commits?apiVersion=2022-11-28
func (se *Engine) SearchCommits(dateQualifier, qFmt string, args ...any) ([]*github.CommitResult, error) {
	query := se.makeQuery(dateQualifier, qFmt, args)
	opts := makeSearchOptions()
	var lst []*github.CommitResult
	for {
		results, resp, err := se.ghClient.Search.Commits(se.ctx, query, opts)
		if err != nil {
			return nil, err
		}
		lst = append(lst, results.Commits...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return lst, nil
}

func (se *Engine) makeQuery(dateQualifier string, qFmt string, args []any) string {
	return fmt.Sprintf(
		"%s:%s..%s %s",
		dateQualifier,
		se.dateRange.StartAsTime().Format(types.DayFormat1),
		se.dateRange.EndAsTime().Format(types.DayFormat1),
		fmt.Sprintf(qFmt, args...))
}

func makeSearchOptions() *github.SearchOptions {
	return &github.SearchOptions{
		Sort:        "",
		Order:       "", // "asc", "desc"
		TextMatch:   false,
		ListOptions: MakeListOptions(),
	}
}

func MakeListOptions() github.ListOptions {
	return github.ListOptions{PerPage: 50}
}
