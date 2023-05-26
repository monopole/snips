package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal"
	"golang.org/x/oauth2"
)

//go:embed README.md
var readMeMd string

const (
	githubDotCom    = "github.com"
	defaultDayCount = 14 // two weeks
)

type pgmArgs struct {
	user     string
	token    string
	dayStart time.Time
	dayCount int
	domain   string
}

func parseArgs() (result pgmArgs) {
	flag.StringVar(&result.domain, "domain", githubDotCom, "the github domain")
	flag.Parse()
	if flag.NArg() < 2 || flag.NArg() > 5 {
		fmt.Print(readMeMd)
		os.Exit(0)
	}
	result.token = flag.Arg(0)
	result.user = flag.Arg(1)
	result.dayStart = time.Now().Round(24 * time.Hour)
	if flag.NArg() > 2 {
		result.dayStart = internal.ParseDate(flag.Arg(2))
	}
	result.dayCount = defaultDayCount
	if flag.NArg() > 3 {
		result.dayCount = internal.ParseDayCount(flag.Arg(3))
	}
	return
}

func main() {
	args := parseArgs()
	ctx := context.Background()
	cl, err := makeClient(ctx, args.domain, args.token)
	if err != nil {
		log.Fatalf("trouble making github client: %s", err.Error())
	}
	questioner{
		user:      args.user,
		dateStart: args.dayStart,
		dateEnd:   args.dayStart.AddDate(0, 0, args.dayCount),
		ctx:       ctx,
		client:    cl,
	}.doIt()
}

func makeClient(ctx context.Context, domain string, token string) (*github.Client, error) {
	oaCl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if domain == githubDotCom {
		return github.NewClient(oaCl), nil
	}
	const scheme = "https://"
	return github.NewEnterpriseClient(scheme+domain, scheme+domain, oaCl)
}
