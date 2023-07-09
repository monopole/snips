package types

import (
	"time"
)

type RepoId struct {
	Org  string
	Name string
}

func (id RepoId) String() string {
	return id.Org + "/" + id.Name
}

func (id RepoId) Equals(other RepoId) bool {
	return id.Org == other.Org && id.Name == other.Name
}

// MyGhOrg is a GitHub Organization.
type MyGhOrg struct {
	Name  string
	Login string
}

// MyIssue holds an issue or a pull request.
// In GitHub, at a high level, an issue and a pull request has the same representation.
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

type IssueSet struct {
	Domain string
	Groups map[RepoId][]MyIssue
}

func (is *IssueSet) Count() int {
	c := 0
	for _, v := range is.Groups {
		c += len(v)
	}
	return c
}

func (is *IssueSet) IsEmpty() bool {
	return is.Count() == 0
}

type MyUser struct {
	Name            string
	Company         string
	Login           string
	Email           string
	GhOrgs          []MyGhOrg
	IssuesCreated   *IssueSet
	IssuesClosed    *IssueSet
	IssuesCommented *IssueSet
	PrsReviewed     *IssueSet
	Commits         map[RepoId][]*MyCommit
}

type Report struct {
	Title      string
	DomainGh   string
	DomainJira string
	Dr         *DayRange
	Users      []*MyUser
}
