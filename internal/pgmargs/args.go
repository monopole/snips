package pgmargs

import (
	"flag"
	"fmt"
	"github.com/monopole/snips/internal/types"
	"os"
)

const (
	GithubPublic              = "github.com"
	githubTesla               = "github.tesla.com"
	githubPublicOAuthClientId = "6019f1c21d0470ec327d"
	githubTeslaOAuthClientId  = "3bfc36851715c5de6d23"

	EnvGhToken = "GH_TOKEN"

	defaultDayCount = 14 // two weeks

	FlagToken    = "token"
	flagDayStart = "day-start"
	flagDayCount = "day-count"
	flagNoEcho   = "suppress-token-echo"
)

// pgmArgs holds clean arguments from the command line.
type pgmArgs struct {
	User        []string
	DateRange   *types.DayRange
	GhDomain    string
	ClientId    string
	Token       string
	NoTokenEcho bool
	CaPath      string
}

// ParseArgs parses and validates arguments from the command line.
func ParseArgs() (*pgmArgs, error) {
	var (
		err      error
		result   pgmArgs
		dayStart string
		dayCount int
	)

	flag.IntVar(&dayCount, flagDayCount, defaultDayCount, "how many days, inclusive of start date")
	flag.StringVar(&dayStart, flagDayStart, "", "the day to start, formatted as "+types.DateOptions())
	flag.StringVar(&result.CaPath, "ca-path", "", "local path to cert file for TLS in oauth dance")
	flag.StringVar(&result.GhDomain, "domain", GithubPublic, "the github domain")
	flag.StringVar(&result.ClientId, "client-id", "", "the oauth clientID")
	flag.StringVar(&result.Token, FlagToken, "",
		fmt.Sprintf("access token for the given GitHub domain (overrides env var %s)", EnvGhToken))

	// By default, the token echoed to stderr, so the user can copy/paste and use with the --token flag
	// to avoid repeated logins.
	flag.BoolVar(&result.NoTokenEcho, flagNoEcho,
		false, "don't echo the access token to stderr for reuse (over the shoulder security)")

	flag.Parse()

	if flag.NArg() == 0 {
		return nil, fmt.Errorf("must specify at least one user")
	}
	// All the arguments should be usernames.
	result.User = flag.Args()

	if result.Token == "" {
		// For fun, see if it can be pulled from env var.
		result.Token = os.Getenv(EnvGhToken)
	}

	if result.ClientId == "" {
		result.ClientId, err = determineClientIdFromDomain(result.GhDomain)
		if err != nil {
			return nil, err
		}
	}

	result.DateRange, err = types.MakeDayRange(dayStart, dayCount)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// determineClientIdFromDomain returns a hardcoded clientId as a function of the domain.
// The only time a user would specify the clientId on the command line would be when
// registering / re-registering this program with some GitHub server.
func determineClientIdFromDomain(domain string) (string, error) {
	if domain == githubTesla {
		return githubTeslaOAuthClientId, nil
	}
	if domain != GithubPublic {
		return "", fmt.Errorf("i have no client id registered with %s", domain)
	}
	return githubPublicOAuthClientId, nil
}

func EchoToken(prefix, token string) {
	fmt.Fprintf(
		os.Stderr, `
%s     In subsequent calls, add this flag:  --%s %s
%s                     or export this var:  export %s=%s
%s To suppress this reminder, use --%s
`[1:],
		prefix, FlagToken, token,
		prefix, EnvGhToken, token,
		prefix, flagNoEcho)
}
