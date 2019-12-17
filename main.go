package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)


const  (
	org1  = "kubernetes-sigs"
	org2  = "monopole"
	repo1 = "kustomize"
	repo2 = "mdrip"
)
// reviews  query:
//  -author:monopole commenter:monopole updated:2019-11-25..2019-12-01
//
// PR query:
//  author:monopole is:pr merged:2019-11-25..2019-12-01

func main() {
	fmt.Println("hey")
	token := os.Getenv("GITHUB_AUTH_TOKEN")
	if token == "" {
		 log.Fatal("Unauthorized: No token present")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	//reportOrgs(client)
	//reportUser(client, ctx)
	reportReviews(client, ctx)
	reportPrs(client, ctx)
}

func reportReviews(client *github.Client, ctx context.Context) {

		opts := &github.SearchOptions{
			Sort:        "",
			Order:       "",
			TextMatch:   false,
			ListOptions: github.ListOptions{
				Page:    0,
				PerPage: 100,
			},
		}
	results, _, err := client.Search.Issues(
		ctx,
		"-author:monopole commenter:monopole updated:2019-11-25..2019-12-01",
		opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("## Reviews")
	for _, r := range results.Issues {
		fmt.Printf(
			" - [%s](https://github.com/kubernetes-sigs/kustomize/pull/%d) \n",
			r.GetTitle(), r.GetNumber())
	}
}
func reportPrs(client *github.Client, ctx context.Context) {
	opts := &github.PullRequestListOptions{
		State:       "closed",
		Head:        "",
		Base:        "",
		Sort:        "",
		Direction:   "desc",
		ListOptions: github.ListOptions{
			Page:    0,
			PerPage: 100,
		},
	}
	results, _, err := client.PullRequests.List(ctx, org1, repo1, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("## PRS")
	for _, r := range results {
		if r.GetUser().GetLogin() == "monopole" {
			t := r.GetMergedAt().Format("2006-01-02")
			fmt.Printf(
				" - [%s](https://github.com/kubernetes-sigs/kustomize/pull/%d) <!-- %s -->\n",
				r.GetTitle(), r.GetNumber(), t)
		}
	}
}

func reportOrgs(client *github.Client) {
	orgs, _, err := client.Organizations.List(
		context.Background(), org2, nil)
	if err != nil {
		log.Fatal(err)
	}

	for i, organization := range orgs {
		fmt.Printf("%v. %v\n", i+1, organization.GetLogin())
	}
}


func reportUser(client *github.Client, ctx context.Context) {
	user, _, err := client.Users.Get(ctx, "monopole")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", user.GetName())
}
