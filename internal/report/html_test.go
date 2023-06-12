package report_test

import (
	"bytes"
	"fmt"
	. "github.com/monopole/snips/internal/report"
	"github.com/monopole/snips/internal/types"
	"testing"
	"time"
)

const (
	urlPr1 = "https://github.tesla.com/design-technology/3dx/pull/636"
	urlPr2 = "https://github.tesla.com/design-technology/argocd-manifests/pull/2555"

	urlCommit1 = "https://github.tesla.com/design-technology/argocd-manifests/pull/2663/commits/fc25519428f4f91813d5a8c324c73ada2d94b578"
	urlCommit2 = "https://github.tesla.com/design-technology/argocd-manifests/pull/2663/commits/bbd9f61f0c1bb26e58641f15da872afce9f6c1ec"

	ts1 = "13 Jun 19 10:11 PST"
	ts2 = "15 Jun 19 10:17 PST"

	orgName1 = "federationOfPlanets"
	orgName2 = "bitCoinLosers"

	title1 = "Fry the older bananas"
	title2 = "Indemnify the cheese eaters"
)

var (
	org1 = types.MyOrg{Name: orgName1, Login: "Micheal"}
	org2 = types.MyOrg{Name: orgName2, Login: "Barton"}

	repoId1 = types.RepoId{
		Org:  orgName1,
		Repo: "marsToilet",
	}
	repoId2 = types.RepoId{
		Org:  orgName2,
		Repo: "jupiterToast",
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
		Pr:               nil,
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

func Test_WriteHtml(t *testing.T) {
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
				Orgs:    []types.MyOrg{org1},
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
			WriteHtml(&b, &types.Report{
				Title: "foosball",
				Dr:    dr,
				Users: []*types.MyUser{&tt.dude},
			})
			fmt.Print(b.String())
			if b.String() != tt.result {
				t.Fail()
			}
		})
	}
}
