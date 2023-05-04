package internal

import "github.com/google/go-github/v52/github"

func MakeListOptions() github.ListOptions {
	return github.ListOptions{PerPage: 50}
}

func MakeSearchOptions() *github.SearchOptions {
	return &github.SearchOptions{
		Sort:        "",
		Order:       "", // "asc", "desc"
		TextMatch:   false,
		ListOptions: MakeListOptions(),
	}
}
