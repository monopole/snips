package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/fake"
	"github.com/monopole/snips/internal/mygh/client"
	"github.com/monopole/snips/internal/mygh/oauth"
	"github.com/monopole/snips/internal/mygh/search"
	"github.com/monopole/snips/internal/myhttp"
	"github.com/monopole/snips/internal/myjira"
	"github.com/monopole/snips/internal/pgmargs"
	"github.com/monopole/snips/internal/report/html"
	"github.com/monopole/snips/internal/report/md"
	"github.com/monopole/snips/internal/types"
)

//go:embed README.md
var readMeMd string

func main() {
	args, err := pgmargs.ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** usage error: %s\n\n", err.Error())
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}
	if !args.TestRenderOnly && !args.JustGetGhToken && len(args.UserNames) == 0 {
		fmt.Fprintf(os.Stderr, readMeMd)
		os.Exit(0)
	}

	var users []*types.MyUser
	if args.TestRenderOnly {
		users = fake.MakeSliceOfFakeUserData()
	} else {
		if users, err = getUserData(args); err != nil {
			log.Fatalf(err.Error())
		}
	}

	writeF := html.WriteHtmlReport
	if args.Markdown {
		writeF = md.WriteMdReport
	}
	if err = writeF(
		os.Stdout,
		&types.Report{
			Title:      args.Title,
			DomainGh:   args.Gh.Domain,
			DomainJira: args.Jira.Domain,
			Dr:         args.DateRange,
			Users:      users,
		}); err != nil {
		log.Fatalf(err.Error())
	}
}

func getUserData(args *pgmargs.Args) ([]*types.MyUser, error) {
	htCl, err := myhttp.MakeHttpClient(args.CaPath)
	if err != nil {
		return nil, err
	}
	if args.JustGetGhToken || args.Gh.Token == "" {
		args.Gh.Token, err = oauth.GetAccessToken(&oauth.Params{
			GhDomain: args.Gh.Domain,
			ClientId: args.Gh.ClientId,
			HttpCl:   htCl,
			Verbose:  false,
		})
		if err != nil {
			return nil, err
		}
		if args.JustGetGhToken {
			fmt.Println(args.Gh.Token)
			return nil, nil
		}
		if !args.NoTokenEcho {
			pgmargs.EchoToken(oauth.WarningPrefix, args.Gh.Token)
		}
	}
	ctx := context.Background()
	var (
		ghCl  *github.Client
		users []*types.MyUser
	)
	ghCl, err = client.MakeGhApiClient(ctx, args.Gh.Domain, args.Gh.Token)
	if err != nil {
		return nil, fmt.Errorf("trouble making github client: %w", err)
	}
	users, err = search.MakeEngine(
		ctx, ghCl, args.Gh.Domain).LookupPeeps(args.UserNames, args.DateRange)
	if err != nil {
		return nil, fmt.Errorf("trouble doing queries: %w", err)
	}

	if args.Jira.Token != "" {
		err = myjira.MakeJiraBoss(
			htCl, &args.Jira, args.DateRange).DoSearch(users)
		if err != nil {
			return nil, err
		}
	}
	return users, nil
}
