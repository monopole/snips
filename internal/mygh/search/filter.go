package search

import (
	"github.com/google/go-github/v52/github"
	"log"
)

type myFilter int

const (
	unknown myFilter = iota
	keepOnlyPrs
	rejectPrs
)

func (f myFilter) from(issues []*github.Issue) (result []*github.Issue) {
	if f != keepOnlyPrs && f != rejectPrs {
		log.Fatalf("unable to deal with filter %v", f)
	}
	for _, issue := range issues {
		if issue.IsPullRequest() == (f == keepOnlyPrs) {
			result = append(result, issue)
		}
	}
	return
}
