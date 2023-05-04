package internal

import (
	"context"
	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

func MakeClient(ctx context.Context, token string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
