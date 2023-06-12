package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/monopole/snips/internal/search"
	"github.com/monopole/snips/internal/types"
	"log"
	"os"

	"github.com/monopole/snips/internal/report"

	"github.com/monopole/snips/internal/query"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/oauth"
	"github.com/monopole/snips/internal/pgmargs"
	"golang.org/x/oauth2"
)

//go:embed README.md
var readMeMd string

func main() {
	args, err := pgmargs.ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** Error: %s\n\n", err.Error())
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}
	if !args.JustGetToken && len(args.User) == 0 {
		fmt.Fprintf(os.Stderr, readMeMd)
		os.Exit(0)
	}
	if args.JustGetToken || args.Token == "" {
		args.Token, err = oauth.GetAccessToken(&oauth.Params{
			GhDomain: args.GhDomain,
			ClientId: args.ClientId,
			CaPath:   args.CaPath,
			Verbose:  false,
		})
		if err != nil {
			log.Fatal(err)
		}
		if args.JustGetToken {
			fmt.Println(args.Token)
			os.Exit(0)
		}
		if !args.NoTokenEcho {
			pgmargs.EchoToken(oauth.WarningPrefix, args.Token)
		}
	}
	if len(args.User) == 0 {
		fmt.Fprintf(os.Stderr, readMeMd)
		os.Exit(1)
	}
	ctx := context.Background()
	cl, err := makeApiClient(ctx, args.GhDomain, args.Token)
	if err != nil {
		log.Fatalf("trouble making github client: %s", err.Error())
	}
	result, err := (&query.Worker{
		Users:    args.User,
		Se:       search.MakeEngine(ctx, cl, args.DateRange),
		Ctx:      ctx,
		GhClient: cl,
	}).DoIt()
	if err != nil {
		log.Fatalf("trouble doing queries: %s", err.Error())
	}
	title := args.Title
	if title == "" {
		title = "Activity at " + args.GhDomain
	}
	if args.Markdown {
		report.PrintReport(&types.Report{
			Title: title,
			Dr:    args.DateRange,
			Users: result,
		})
	} else {
		err = report.WriteHtmlReport(os.Stdout, &types.Report{
			Title: title,
			Dr:    args.DateRange,
			Users: result,
		})
		if err != nil {
			log.Fatalf("trouble rendering html: %s", err.Error())
		}
	}
}

func makeApiClient(ctx context.Context, domain string, token string) (*github.Client, error) {
	oaCl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if domain == pgmargs.GithubPublic {
		return github.NewClient(oaCl), nil
	}
	const scheme = "https://"
	return github.NewEnterpriseClient(scheme+domain, scheme+domain, oaCl)
}
