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
	githubPublic              = "github.com"
	githubTesla               = "github.tesla.com"
	githubPublicOAuthClientId = "6019f1c21d0470ec327d"
	githubTeslaOAuthClientId  = "3bfc36851715c5de6d23"

	defaultDayCount = 14 // two weeks
	flagToken       = "token"
)

type pgmArgs struct {
	user        string
	dayStart    time.Time
	dayCount    int
	domain      string
	clientId    string
	token       string
	noTokenEcho bool
}

func parseArgs() (result pgmArgs) {
	flag.StringVar(&result.domain, "domain", githubPublic, "the github domain")
	flag.StringVar(&result.clientId, "client-id", "", "the oauth clientID")
	flag.StringVar(&result.token, flagToken, "", "access token for the given github domain")
	flag.BoolVar(&result.noTokenEcho, "suppress-token-echo", false, "don't echo the access token to stdout for reuse")
	flag.Parse()
	if flag.NArg() < 1 || flag.NArg() > 3 {
		fmt.Print(readMeMd)
		os.Exit(0)
	}
	result.user = flag.Arg(0)
	result.dayStart = time.Now().Round(24 * time.Hour)
	if flag.NArg() > 1 {
		result.dayStart = internal.ParseDate(flag.Arg(1))
	}
	result.dayCount = defaultDayCount
	if flag.NArg() > 2 {
		result.dayCount = internal.ParseDayCount(flag.Arg(2))
	}
	return
}

func guessClientId(override, domain string) string {
	if override != "" {
		return override
	}
	if domain == githubTesla {
		return githubTeslaOAuthClientId
	}
	if domain != githubPublic {
		log.Fatalf("i have no client id registered with %s", domain)
	}
	return githubPublicOAuthClientId
}

func main() {
	var err error
	args := parseArgs()
	if args.token == "" {
		args.token, err = internal.GetOAuthAccessToken(&internal.OAuthParams{
			Domain:   args.domain,
			ClientId: guessClientId(args.clientId, args.domain),
		})
		if err != nil {
			log.Fatal(err)
		}
		if !args.noTokenEcho {
			fmt.Fprintf(
				os.Stderr, "  ***** In subsequent calls, add the flag:  --%s %s\n\n\n", flagToken, args.token)
		}
	}
	ctx := context.Background()
	cl, err := makeApiClient(ctx, args.domain, args.token)
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

func makeApiClient(ctx context.Context, domain string, token string) (*github.Client, error) {
	oaCl := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	if domain == githubPublic {
		return github.NewClient(oaCl), nil
	}
	const scheme = "https://"
	return github.NewEnterpriseClient(scheme+domain, scheme+domain, oaCl)
}
