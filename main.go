package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v28/github"
	"golang.org/x/oauth2"
)

const user = "monopole" // "monopole"
// reviews  query:
//  author:monopole commenter:monopole updated:2019-11-25..2019-12-01
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
	reportOrgs(client)
	reportUser(client, ctx)
reportPrs()
	commitResult, _, err := client.Search.Commits(
		ctx,
		"author:monopole is:pr merged:2019-11-25..2019-12-01",
		&github.SearchOptions{
			Sort:        "",
			Order:       "",
			TextMatch:   false,
			ListOptions: github.ListOptions{},
		})
	if err != nil {
		log.Fatal(err)
	}
	for i, r := range commitResult.Commits {
		fmt.Printf("%v. %v\n", i+1, r.GetCommit().GetMessage())
	}
}

func reportOrgs(client *github.Client) {
	orgs, _, err := client.Organizations.List(
		context.Background(), user, nil)
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
	fmt.Printf("%v\n", user)
}

func reportUser(client *github.Client, ctx context.Context) {
	user, _, err := client.Users.Get(ctx, "monopole")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v\n", user)
}
