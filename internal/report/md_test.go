package report_test

import (
	"bytes"
	. "github.com/monopole/snips/internal/report"
	"github.com/monopole/snips/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WriteMdIssue(t *testing.T) {
	tests := map[string]struct {
		issue  types.MyIssue
		result string
	}{
		"t1": {
			issue:  issue1,
			result: "`2019-Jun-13` [Fry the older bananas](https://github.tesla.com/design-technology/3dx/pull/636)",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteMdIssue(&b, &tt.issue)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteMdCommit(t *testing.T) {
	tests := map[string]struct {
		commit types.MyCommit
		result string
	}{
		"t1": {
			commit: commit1,
			result: "`2019-Jun-13` [`fc25519`](https://github.tesla.com/design-technology/argocd-manifests/pull/2663/commits/fc25519428f4f91813d5a8c324c73ada2d94b578) (pull/[600](https://github.tesla.com/design-technology/3dx/pull/636)) Fry the older bananas",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteMdCommit(&b, &tt.commit)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteMdLabelledIssueMap(t *testing.T) {
	tests := map[string]struct {
		l      string
		m      map[types.RepoId][]types.MyIssue
		result string
	}{
		"t1": {
			l:      "issues commented",
			result: `### No issues commented`,
		},
		"t2": {
			l: "issues reviewed",
			m: map[types.RepoId][]types.MyIssue{
				repoId1: {issue1, issue2},
				repoId2: {issue1, issue2},
			},
			result: `### issues reviewed:

#### bitCoinLosers/jupiterToast

  - ` + "`2019-Jun-13`" + ` [Fry the older bananas](https://github.tesla.com/design-technology/3dx/pull/636)
  - ` + "`2019-Jun-15`" + ` [Indemnify the cheese eaters](https://github.tesla.com/design-technology/argocd-manifests/pull/2555)

#### federationOfPlanets/marsToilet

  - ` + "`2019-Jun-13`" + ` [Fry the older bananas](https://github.tesla.com/design-technology/3dx/pull/636)
  - ` + "`2019-Jun-15`" + ` [Indemnify the cheese eaters](https://github.tesla.com/design-technology/argocd-manifests/pull/2555)
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteMdLabelledIssueMap(&b, tt.l, tt.m)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteMdLabelledCommitMap(t *testing.T) {
	tests := map[string]struct {
		l      string
		m      map[types.RepoId][]*types.MyCommit
		result string
	}{
		"t1": {
			l:      "commits",
			result: `### No commits`,
		},
		"t2": {
			l: "commits",
			m: map[types.RepoId][]*types.MyCommit{
				repoId1: {&commit1, &commit2},
			},
			result: `### commits

#### federationOfPlanets/marsToilet

` + " - `2019-Jun-13` [`fc25519`]" + `(https://github.tesla.com/design-technology/argocd-manifests/pull/2663/commits/fc25519428f4f91813d5a8c324c73ada2d94b578) (pull/[600](https://github.tesla.com/design-technology/3dx/pull/636)) Fry the older bananas
` + " - `2019-Jun-13` [`bbd9f61`]" + `(https://github.tesla.com/design-technology/argocd-manifests/pull/2663/commits/bbd9f61f0c1bb26e58641f15da872afce9f6c1ec) Fry the older bananas
`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteMdLabelledCommitMap(&b, tt.l, tt.m)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteMdReport(t *testing.T) {
	tests := map[string]struct {
		dude   types.MyUser
		result string
	}{
		"t1": {
			dude: types.MyUser{
				Name:    "Bobby Bobface",
				Company: "TESLA",
				Login:   "bobby",
				Email:   "bob@tesla.com",
				Orgs:    []types.MyOrg{org1, org2},
				IssuesCreated: map[types.RepoId][]types.MyIssue{
					repoId1: {issue1, issue2},
					repoId2: {issue1, issue2},
				},
				IssuesClosed:    nil,
				IssuesCommented: nil,
				PrsReviewed:     nil,
				Commits: map[types.RepoId][]*types.MyCommit{
					repoId1: {&commit1, &commit2},
				},
			},
			result: "hey there",
		},
	}
	dr, err := types.MakeDayRange("2023/01/03", "", 0)
	if err != nil {
		t.Fatalf("bad time: %s", err.Error())
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			assert.NoError(t, WriteMdReport(&b, &types.Report{
				Title:    "hello I am the report title",
				GhDomain: "github.com",
				Dr:       dr,
				Users:    []*types.MyUser{&tt.dude},
			}))
			//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++")
			//fmt.Println(b.String())
			//fmt.Println("+++++++++++++++++++++++++++++++++++++++++++++++++")
			//assert.Equal(t, tt.result, b.String())
		})
	}
}
