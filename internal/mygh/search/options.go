package search

import "github.com/google/go-github/v52/github"

func makeSearchOptions() *github.SearchOptions {
	return &github.SearchOptions{
		Sort:        "",
		Order:       "", // "asc", "desc"
		TextMatch:   false,
		ListOptions: makeListOptions(),
	}
}

func makeListOptions() github.ListOptions {
	return github.ListOptions{PerPage: 50}
}
