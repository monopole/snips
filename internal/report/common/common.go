package common

import (
	"strings"
	"time"

	"github.com/monopole/snips/internal/types"
)

// MakeFuncMap makes a string to function map for use in Go template rendering.
func MakeFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"toUpper": strings.ToUpper,
		"shaSmall": func(s string) string {
			return s[0:7]
		},
		"snipDate": func(t time.Time) string {
			return t.Format(types.DayFormatHuman)
		},
		"prettyDateRange": func(dr *types.DayRange) string {
			return dr.PrettyRange()
		},
		"labeledIssueSet":  LabeledIssueSet,
		"labeledCommitMap": LabeledCommitMap,
		"mapTotalCommits": func(m map[types.RepoId][]*types.MyCommit) int {
			c := 0
			for _, v := range m {
				c += len(v)
			}
			return c
		},
		"bigEnough": func(s int) bool {
			return s > 5
		},
		"domainsAndUser": func(dGh string, dJira string, u *types.MyUser) interface{} {
			return &struct {
				Dgh   string
				Djira string
				U     *types.MyUser
			}{Dgh: dGh, Djira: dJira, U: u}
		},
		"countAndItemName": func(c int, n string) interface{} {
			return &struct {
				C int
				N string
			}{C: c, N: n}
		},
		"domainAndCommitMap": func(dGh string, m map[types.RepoId][]*types.MyCommit) interface{} {
			return &struct {
				Dgh string
				M   map[types.RepoId][]*types.MyCommit
			}{Dgh: dGh, M: m}
		},
		"domainAndRepo": func(dGh string, rid types.RepoId) interface{} {
			return &DomainAndRepo{
				Dgh: dGh,
				Rid: rid,
			}
		},
		"domainAndOrgs": func(dGh string, o []types.MyGhOrg) interface{} {
			return &struct {
				Dgh    string
				GhOrgs []types.MyGhOrg
			}{Dgh: dGh, GhOrgs: o}
		},
		"lowerHyphen": func(what string) string {
			return strings.ReplaceAll(strings.ToLower(what), " ", "-")
		},
	}
}

type DomainAndRepo struct {
	Dgh string
	Rid types.RepoId
}

func (dr DomainAndRepo) HRef() string {
	if strings.Contains(dr.Dgh, "github") {
		// Try to make a GitHub link.
		return dr.Dgh + "/" + dr.Rid.String()
	}
	// Try to make a Jira link.
	return dr.Dgh + "/projects/" + dr.Rid.Name + "/issues"
}

func LabeledCommitMap(l string, dGh string, m map[types.RepoId][]*types.MyCommit) interface{} {
	return &struct {
		Label string
		Dgh   string
		M     map[types.RepoId][]*types.MyCommit
	}{Label: l, Dgh: dGh, M: m}
}

func LabeledIssueSet(l string, iSet *types.IssueSet) interface{} {
	return &struct {
		Label string
		ISet  *types.IssueSet
	}{Label: l, ISet: iSet}
}
