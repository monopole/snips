package internal

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

func ParseDate(arg string) time.Time {
	t, err := time.Parse("2006-01-02", arg)
	if err != nil {
		fmt.Printf("Trouble with date specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return t
}

func ParseDayCount(arg string) int {
	i, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Printf("Trouble with day count specification: '%s'\n", arg)
		log.Fatal(err)
	}
	return i
}

func MakeClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
