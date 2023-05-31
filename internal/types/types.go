package types

import (
	"time"
)

type RepoName string

type MyOrg struct {
	Name  string
	Login string
}

// MyIssue holds issues, pull requests and commits.
// In the case of commits, Title holds the first line of the commit message,
// Updated holds the commit date, and Number is left empty.
// TODO: In the commits section of the report, include only commits that lack a PR.
// In the PR section, add links to commits if there's more than one commit in the PR.
type MyIssue struct {
	Number  int
	Title   string
	HtmlUrl string
	Updated time.Time
}

type MyUser struct {
	Name            string
	Company         string
	Login           string
	Email           string
	Orgs            []MyOrg
	IssuesCreated   map[RepoName][]MyIssue
	IssuesClosed    map[RepoName][]MyIssue
	IssuesCommented map[RepoName][]MyIssue
	PrsMerged       map[RepoName][]MyIssue
	PrsReviewed     map[RepoName][]MyIssue
	Commits         map[RepoName][]MyIssue
}
