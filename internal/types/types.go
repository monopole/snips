package types

import (
	"time"
)

type RepoId struct {
	Org  string
	Repo string
}

func (id RepoId) String() string {
	return id.Org + "/" + id.Repo
}

func (id RepoId) Equals(other RepoId) bool {
	return id.Org == other.Org && id.Repo == other.Repo
}

type MyOrg struct {
	Name  string
	Login string
}

// MyIssue holds an issue or a pull request, since in GitHub those are the same at a high level.
type MyIssue struct {
	RepoId  RepoId
	Number  int
	Title   string
	HtmlUrl string
	Updated time.Time
}

type MyCommit struct {
	RepoId           RepoId
	Sha              string
	Url              string
	MessageFirstLine string
	Committed        time.Time
	Author           string
	Pr               *MyIssue
}

type MyUser struct {
	Name            string
	Company         string
	Login           string
	Email           string
	Orgs            []MyOrg
	IssuesCreated   map[RepoId][]MyIssue
	IssuesClosed    map[RepoId][]MyIssue
	IssuesCommented map[RepoId][]MyIssue
	PrsReviewed     map[RepoId][]MyIssue
	Commits         map[RepoId][]*MyCommit
}

type Report struct {
	Title string
	Dr    *DayRange
	Users []*MyUser
}
