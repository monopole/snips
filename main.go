package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/mygh/client"
	"github.com/monopole/snips/internal/mygh/oauth"
	"github.com/monopole/snips/internal/mygh/search"
	"github.com/monopole/snips/internal/myhttp"
	"github.com/monopole/snips/internal/myjira"
	"github.com/monopole/snips/internal/pgmargs"
	"github.com/monopole/snips/internal/report"
	"github.com/monopole/snips/internal/types"
)

//go:embed README.md
var readMeMd string

const (
	doTheJira = true
)

func trueMain(args *pgmargs.Args) error {
	htCl, err := myhttp.MakeHttpClient(args.CaPath)
	if err != nil {
		return err
	}
	ctx := context.Background()

	if args.JustGetGhToken || args.Gh.Token == "" {
		args.Gh.Token, err = oauth.GetAccessToken(&oauth.Params{
			GhDomain: args.Gh.Domain,
			ClientId: args.Gh.ClientId,
			HttpCl:   htCl,
			Verbose:  false,
		})
		if err != nil {
			return err
		}
		if args.JustGetGhToken {
			fmt.Println(args.Gh.Token)
			return nil
		}
		if !args.NoTokenEcho {
			pgmargs.EchoToken(oauth.WarningPrefix, args.Gh.Token)
		}
	}
	if len(args.User) == 0 {
		return fmt.Errorf("no users specified")
	}
	if doTheJira && args.Jira.Token == "" {
		return fmt.Errorf(pgmargs.HelpJiraToken)
	}
	var (
		ghCl  *github.Client
		users []*types.MyUser
	)
	ghCl, err = client.MakeGhApiClient(ctx, args.Gh.Domain, args.Gh.Token)
	if err != nil {
		return fmt.Errorf("trouble making github client: %w", err)
	}
	users, err = search.MakeEngine(ctx, ghCl).LookupPeeps(args.User, args.DateRange)
	if err != nil {
		return fmt.Errorf("trouble doing queries: %w", err)
	}

	if doTheJira {
		jb := myjira.MakeJiraBoss(htCl, args)
		err = jb.DoSearch(users)
		if err != nil {
			return err
		}
	}

	writeF := report.WriteHtmlReport
	if args.Markdown {
		writeF = report.WriteMdReport
	}
	if err = writeF(
		os.Stdout,
		&types.Report{
			Title:    args.Title,
			GhDomain: args.Gh.Domain,
			Dr:       args.DateRange,
			Users:    users,
		}); err != nil {
		return fmt.Errorf("trouble rendering report: %w", err)
	}

	return nil
}

func main() {
	args, err := pgmargs.ParseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "*** Error: %s\n\n", err.Error())
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}
	if !args.JustGetGhToken && len(args.User) == 0 {
		fmt.Fprintf(os.Stderr, readMeMd)
		os.Exit(0)
	}
	if err = trueMain(args); err != nil {
		log.Fatalf(err.Error())
	}
}
