package types

import (
	"time"
)

type RepoName string

type MyOrg struct {
	Name  string
	Login string
}

type MyIssue struct {
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
}
