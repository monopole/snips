package internal

// Not using this stuff, but want to be sure it compiles.

import (
	"context"
	"fmt"
	"github.com/google/go-github/v52/github"
	"log"
	"time"
)

type oldQuestioner struct {
	user      string
	dateStart time.Time
	dateEnd   time.Time
	client    *github.Client
	ctx       context.Context
}

type ghRepo struct {
	org  string
	name string
}

// Should get this from command line if we use it, but at the moment
// not using an API that needs the repo.  Instead, we ask github
// for all activity of a given user in a given date range.
var myRepos = []ghRepo{
	{"kubernetes-sigs", "kustomize"},
	{"kubernetes-sigs", "cli-utils"},
	{"GoogleContainerTools", "kpt"},
}

// Not using this any more
func (q *oldQuestioner) reportPrs() {
	fmt.Print("\n## PRS\n\n")
	for _, repo := range myRepos {
		prs := q.queryPrs(repo)
		if len(prs) > 0 {
			//prList = sortPrsDate(prList)
			fmt.Printf("#### %s\n\n", repo.name)
			for _, pr := range prs {
				PrintPrLink(pr)
			}
			fmt.Println()
		}
	}
}

func (q *oldQuestioner) queryPrs(repo ghRepo) (prs []*github.PullRequest) {
	lOpts := MakeListOptions()
	for {
		prList, resp, err := q.client.PullRequests.List(
			q.ctx,
			repo.org,
			repo.name,
			&github.PullRequestListOptions{
				State:       "closed",
				Head:        q.user + ":master",
				Base:        "",
				Sort:        "created",
				Direction:   "desc",
				ListOptions: lOpts,
			})
		if err != nil {
			log.Fatal(err)
		}
		prs = append(prs, prList...)
		if resp.NextPage == 0 {
			break
		}
		fmt.Print(".")
		lOpts.Page = resp.NextPage
	}
	fmt.Println()
	return
}
