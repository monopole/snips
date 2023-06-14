package pgmargs

import (
	"flag"
	"fmt"
	"os"

	"github.com/monopole/snips/internal/types"
)

const (
	GithubPublic              = "github.com"
	githubTesla               = "github.tesla.com"
	githubPublicOAuthClientId = "6019f1c21d0470ec327d"
	githubTeslaOAuthClientId  = "3bfc36851715c5de6d23"

	envGhToken = "GH_TOKEN"
	flagToken  = "gh-token"

	flagDayStart = "day-start"
	flagDayEnd   = "day-end"
	flagDayCount = "day-count"
	flagNoEcho   = "suppress-token-echo"
	flagMarkdown = "md"
)

// pgmArgs holds clean arguments from the command line.
type pgmArgs struct {
	User         []string
	Title        string
	DateRange    *types.DayRange
	GhDomain     string
	ClientId     string
	Token        string
	NoTokenEcho  bool
	Markdown     bool
	JustGetToken bool
	CaPath       string
}

// ParseArgs parses and validates arguments from the command line.
func ParseArgs() (*pgmArgs, error) {
	var (
		err      error
		result   pgmArgs
		dayStart string
		dayEnd   string
		dayCount int
	)

	flag.IntVar(&dayCount, flagDayCount, 0, "how many days, inclusive of start date")
	flag.StringVar(&dayStart, flagDayStart, "", "the day to start, formatted as "+types.DateOptions())
	flag.StringVar(&dayEnd, flagDayEnd, "", "the day to end, formatted as "+types.DateOptions()+", (default today)")
	flag.StringVar(&result.Title, "title", "", "the title of the report")
	flag.BoolVar(&result.JustGetToken, "get-gh-token", false, "force login, return the gh-token")
	flag.BoolVar(&result.Markdown, flagMarkdown, false, "emit markdown instead of HTML")
	flag.StringVar(&result.CaPath, "ca-path", "", "local path to cert file for TLS in oauth dance")
	flag.StringVar(&result.GhDomain, "domain", GithubPublic, "the github domain")
	flag.StringVar(&result.ClientId, "client-id", "", "the oauth clientID")
	flag.StringVar(&result.Token, flagToken, "",
		fmt.Sprintf("access token for the given GitHub domain (overrides env var %s)", envGhToken))

	// By default, the GitHub access token echoed to stderr, so that the
	// user can copy/paste it into flagToken or envGhToken to avoid repeated logins.
	flag.BoolVar(&result.NoTokenEcho, flagNoEcho,
		false, fmt.Sprintf("don't echo the value of %s (over-the-shoulder security)", flagToken))

	flag.Parse()

	// All the arguments should be usernames.
	result.User = flag.Args()

	if result.Token == "" {
		// For fun, see if it can be pulled from env var.
		result.Token = os.Getenv(envGhToken)
	}

	if result.ClientId == "" {
		result.ClientId, err = determineClientIdFromDomain(result.GhDomain)
		if err != nil {
			return nil, err
		}
	}

	if dayStart != "" && dayEnd != "" && dayCount > 0 {
		return nil, fmt.Errorf("specify any two of --%s, --%s and --%s", flagDayStart, flagDayEnd, flagDayCount)
	}
	result.DateRange, err = types.MakeDayRange(dayStart, dayEnd, dayCount)
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
%s Login successful.
%s     In subsequent calls, add this flag:  --%s %s
%s                     or export this var:  export %s=%s
%s To suppress this reminder, use --%s
`[1:],
		prefix,
		prefix, flagToken, token,
		prefix, envGhToken, token,
		prefix, flagNoEcho)
}
