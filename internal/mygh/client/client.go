package client

import (
	"context"
	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/pgmargs"
	"golang.org/x/oauth2"
)

func MakeGhApiClient(ctx context.Context, domain string, token string) (*github.Client, error) {
	oaCl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if domain == pgmargs.GithubPublic {
		return github.NewClient(oaCl), nil
	}
	const scheme = "https://"
	return github.NewEnterpriseClient(scheme+domain, scheme+domain, oaCl)
}
