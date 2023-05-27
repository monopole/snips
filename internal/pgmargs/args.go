package pgmargs

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
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
)

const (
	dayFormat1 = "2006-01-02"
	dayFormat2 = "2006-JAN-02"
	dayFormat3 = "2006-jan-02"
)

func allDateFormats() []string {
	return []string{dayFormat1, dayFormat2, dayFormat3}
}

// pgmArgs holds clean arguments from the command line.
type pgmArgs struct {
	User        []string
	DayStart    time.Time
	DayCount    int
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
	)

	flag.IntVar(&result.DayCount, flagDayCount, defaultDayCount, "how many days, inclusive of start date")
	flag.StringVar(&dayStart, flagDayStart, "", "the day to start, formatted as "+dateOptions())
	flag.StringVar(&result.CaPath, "ca-path", "", "local path to cert file for TLS in oauth dance")
	flag.StringVar(&result.GhDomain, "domain", GithubPublic, "the github domain")
	flag.StringVar(&result.ClientId, "client-id", "", "the oauth clientID")
	flag.StringVar(&result.Token, FlagToken, "",
		fmt.Sprintf("access token for the given GitHub domain (overrides env var %s)", EnvGhToken))

	// By default, the token echoed to stderr, so the user can copy/paste and use with the --token flag
	// to avoid repeated logins.
	flag.BoolVar(&result.NoTokenEcho, "suppress-token-echo",
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

	if result.DayCount < 1 {
		return nil, fmt.Errorf("--%s must be at least 1", flagDayCount)
	}

	if dayStart == "" {
		// Default is today minus DayCount, revealing recent activity by default.
		result.DayStart = time.Now().Round(24*time.Hour).AddDate(0, 0, -result.DayCount)
	} else {
		result.DayStart, err = parseDate(dayStart)
		if err != nil {
			return nil, fmt.Errorf("bad value in flag --%s %s", flagDayStart, dayStart)
		}
	}

	return &result, nil
}

func dateOptions() string {
	opts := allDateFormats()
	return strings.Join(opts[0:len(opts)-1], ", ") + " or " + opts[len(opts)-1]
}

func parseDate(v string) (time.Time, error) {
	for _, f := range allDateFormats() {
		if t, err := time.Parse(f, v); err == nil {
			return t, nil
		}
	}
	return time.Now(), fmt.Errorf("bad value in flag --%s %s", flagDayStart, v)
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
