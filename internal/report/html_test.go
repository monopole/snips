package report_test

import (
	"bytes"
	"testing"
	"time"

	. "github.com/monopole/snips/internal/report"
	"github.com/monopole/snips/internal/types"
	"github.com/stretchr/testify/assert"
)

const (
	urlPr1 = "https://github.acmecorp.com/design-technology/3dx/pull/636"
	urlPr2 = "https://github.acmecorp.com/design-technology/argocd-manifests/pull/2555"

	urlCommit1 = "https://github.acmecorp.com/design-technology/argocd-manifests/pull/2663/commits/fc25519428f4f91813d5a8c324c73ada2d94b578"
	urlCommit2 = "https://github.acmecorp.com/design-technology/argocd-manifests/pull/2663/commits/bbd9f61f0c1bb26e58641f15da872afce9f6c1ec"

	ts1 = "13 Jun 19 10:11 PST"
	ts2 = "15 Jun 19 10:17 PST"

	orgName1 = "federationOfPlanets"
	orgName2 = "bitCoinLosers"

	title1 = "Fry the older bananas"
	title2 = "Indemnify the cheese eaters"
)

var (
	org1 = types.MyGhOrg{Name: orgName1, Login: "Micheal"}
	org2 = types.MyGhOrg{Name: orgName2, Login: "Barton"}

	repoId1 = types.RepoId{
		Org:  orgName1,
		Name: "marsToilet",
	}
	repoId2 = types.RepoId{
		Org:  orgName2,
		Name: "jupiterToast",
	}

	time1, _ = time.Parse(time.RFC822, ts1)
	time2, _ = time.Parse(time.RFC822, ts2)
	issue1   = types.MyIssue{
		RepoId:  repoId1,
		Number:  600,
		Title:   title1,
		HtmlUrl: urlPr1,
		Updated: time1,
	}
	issue2 = types.MyIssue{
		RepoId:  repoId1,
		Number:  600,
		Title:   title2,
		HtmlUrl: urlPr2,
		Updated: time2,
	}
	commit1 = types.MyCommit{
		RepoId:           repoId1,
		Sha:              "fc25519",
		Url:              urlCommit1,
		MessageFirstLine: "Fry the older bananas",
		Committed:        time1,
		Author:           "bob",
		Pr:               &issue1,
	}
	commit2 = types.MyCommit{
		RepoId:           repoId1,
		Sha:              "bbd9f61",
		Url:              urlCommit2,
		MessageFirstLine: "Fry the older bananas",
		Committed:        time1,
		Author:           "bob",
		Pr:               nil,
	}
)

func Test_WriteHtmlIssue(t *testing.T) {
	tests := map[string]struct {
		issue  types.MyIssue
		result string
	}{
		"t1": {
			issue:  issue1,
			result: "<code>2019-Jun-13</code> &nbsp; <a href=\"https://github.acmecorp.com/design-technology/3dx/pull/636\"> Fry the older bananas </a>",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteHtmlIssue(&b, &tt.issue)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteHtmlCommit(t *testing.T) {
	tests := map[string]struct {
		commit types.MyCommit
		result string
	}{
		"t1": {
			commit: commit1,
			result: `<code>2019-Jun-13
<a href="https://github.acmecorp.com/design-technology/argocd-manifests/pull/2663/commits/fc25519428f4f91813d5a8c324c73ada2d94b578">fc25519</a> (pull/<a href="https://github.acmecorp.com/design-technology/3dx/pull/636">600</a>)
</code>
&nbsp; Fry the older bananas`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteHtmlCommit(&b, &tt.commit)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteHtmlLabeledIssueMap(t *testing.T) {
	tests := map[string]struct {
		l      string
		iSet   *types.IssueSet
		result string
	}{
		"t1": {
			l:      "issues commented",
			result: `<h3> No issues commented </h3>`,
		},
		"t2": {
			l: "issues reviewed",
			iSet: &types.IssueSet{
				Domain: "github.bob.com",
				Groups: map[types.RepoId][]types.MyIssue{
					repoId1: {issue1, issue2},
					repoId2: {issue1, issue2},
				},
			},
			result: `<h3> issues reviewed
<span class="itemCount">(4 issues in 2 repos)</span>
</h3>
<div class="issueMap">
<h4> <a href="https://github.bob.com/bitCoinLosers/jupiterToast"> bitCoinLosers/jupiterToast </a> 
<span class="itemCount">(2 issues)</span>
</h4>

<div class="oneIssue"> <code>2019-Jun-13</code> &nbsp; <a href="https://github.acmecorp.com/design-technology/3dx/pull/636"> Fry the older bananas </a> </div>
<div class="oneIssue"> <code>2019-Jun-15</code> &nbsp; <a href="https://github.acmecorp.com/design-technology/argocd-manifests/pull/2555"> Indemnify the cheese eaters </a> </div><h4> <a href="https://github.bob.com/federationOfPlanets/marsToilet"> federationOfPlanets/marsToilet </a>  </h4>

<div class="oneIssue"> <code>2019-Jun-13</code> &nbsp; <a href="https://github.acmecorp.com/design-technology/3dx/pull/636"> Fry the older bananas </a> </div>
<div class="oneIssue"> <code>2019-Jun-15</code> &nbsp; <a href="https://github.acmecorp.com/design-technology/argocd-manifests/pull/2555"> Indemnify the cheese eaters </a> </div>
</div>`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteHtmlLabeledIssueSet(&b, tt.l, tt.iSet)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteHtmlLabeledCommitMap(t *testing.T) {
	tests := map[string]struct {
		l      string
		m      map[types.RepoId][]*types.MyCommit
		result string
	}{
		"t1": {
			l:      "commits",
			result: `<h3> No commits </h3>`,
		},
		"t2": {
			l: "commits",
			m: map[types.RepoId][]*types.MyCommit{
				repoId1: {&commit1, &commit2},
			},
			result: `<h3> commits  </h3>
<div class="issueMap">
<h4> <a href="https://hoser.github.com/federationOfPlanets/marsToilet"> federationOfPlanets/marsToilet </a> </h4>

<div class="oneIssue"> <code>2019-Jun-13
<a href="https://github.acmecorp.com/design-technology/argocd-manifests/pull/2663/commits/fc25519428f4f91813d5a8c324c73ada2d94b578">fc25519</a> (pull/<a href="https://github.acmecorp.com/design-technology/3dx/pull/636">600</a>)
</code>
&nbsp; Fry the older bananas </div>
<div class="oneIssue"> <code>2019-Jun-13
<a href="https://github.acmecorp.com/design-technology/argocd-manifests/pull/2663/commits/bbd9f61f0c1bb26e58641f15da872afce9f6c1ec">bbd9f61</a>
</code>
&nbsp; Fry the older bananas </div>
</div>`,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var b bytes.Buffer
			WriteHtmlLabeledCommitMap(&b, tt.l, tt.m)
			assert.Equal(t, tt.result, b.String())
		})
	}
}

func Test_WriteHtmlReport(t *testing.T) {
	tests := map[string]struct {
		dude   types.MyUser
		result string
	}{
		"t1": {
			dude: types.MyUser{
				Name:    "Bobby McBobface",
				Company: "ACME CORP",
				Login:   "bobby",
				Email:   "bob@acmecorp.com",
				GhOrgs:  []types.MyGhOrg{org1, org2},
				IssuesCreated: &types.IssueSet{
					Domain: "hoser",
					Groups: map[types.RepoId][]types.MyIssue{
						repoId1: {issue1, issue2},
						repoId2: {issue1, issue2},
					},
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
			assert.NoError(t, WriteHtmlReport(&b, &types.Report{
				Title: "hello I am the report title",
				Dr:    dr,
				Users: []*types.MyUser{&tt.dude},
			}))
			//fmt.Println("-------------------")
			//fmt.Println(b.String())
			//fmt.Println("-------------------")
			//assert.Equal(t, tt.result, b.String())
		})
	}
}
