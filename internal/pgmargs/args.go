package pgmargs

import (
	"flag"
	"fmt"
	"os"

	"github.com/monopole/snips/internal/types"
)

const (
	flagDayStart = "day-start"
	flagDayEnd   = "day-end"
	flagDayCount = "day-count"
	flagNoEcho   = "suppress-token-echo"
	flagMarkdown = "md"

	GithubPublic              = "github.com"
	githubTesla               = "github.tesla.com"
	githubPublicOAuthClientId = "6019f1c21d0470ec327d"
	githubTeslaOAuthClientId  = "3bfc36851715c5de6d23"
	jiraTesla                 = "issues.teslamotors.com"

	envGhToken  = "GH_TOKEN"
	flagGhToken = "gh-token"

	envJiraToken  = "JIRA_TOKEN"
	flagJiraToken = "jira-token"
	HelpJiraToken = "Set " + envJiraToken + " to a personal access token value obtained from https://" + jiraTesla + "/secure/ViewProfile.jspa"
)

// githubArgs holds information needed to contact GitHub (public or enterprise instance).
type githubArgs struct {
	Domain   string
	ClientId string
	Token    string
}

// JiraArgs holds information needed to contact jira (enterprise instance).
type JiraArgs struct {
	Domain   string
	ClientId string
	Token    string
}

// Args holds clean arguments from the command line.
type Args struct {
	// User is a slice of usernames to include in the given report.
	User []string
	// Title is the title for the given report.
	Title     string
	DateRange *types.DayRange
	CaPath    string
	Gh        githubArgs
	Jira      JiraArgs
	// NoTokenEcho if true suppresses echo of the value of a newly discovered GH token.
	NoTokenEcho bool
	// JustGetGhToken allows execution to get a token if no usernames are specified.
	// Further, the output is ONLY the token.
	JustGetGhToken bool
	// Markdown means emit markdown rather than HTML in the report.
	Markdown bool
}

// ParseArgs parses and validates arguments from the command line.
func ParseArgs() (*Args, error) {
	var (
		err      error
		result   Args
		dayStart string
		dayEnd   string
		dayCount int
	)

	flag.IntVar(&dayCount, flagDayCount, 0, "how many days, inclusive of start date")
	flag.StringVar(&dayStart, flagDayStart, "", "the day to start, formatted as "+types.DateOptions())
	flag.StringVar(&dayEnd, flagDayEnd, "", "the day to end, formatted as "+types.DateOptions()+", (default today)")
	flag.StringVar(&result.Title, "title", "", "the title of the report")
	flag.BoolVar(&result.Markdown, flagMarkdown, false, "emit markdown instead of HTML")
	flag.StringVar(&result.CaPath, "ca-path", "", "local path to cert file for TLS in oauth dance")

	flag.BoolVar(&result.JustGetGhToken, "get-gh-token", false, "force github login, return the gh-token")
	flag.StringVar(&result.Gh.Domain, "gh-domain", GithubPublic, "the github domain")
	flag.StringVar(&result.Gh.ClientId, "gh-client-id", "", "the oauth clientID from github")
	flag.StringVar(&result.Gh.Token, flagGhToken, "",
		fmt.Sprintf("access token for the given GitHub domain (overrides env var %s)", envGhToken))

	flag.StringVar(&result.Jira.Domain, "jira-domain", jiraTesla, "the jira domain")
	// To make a jira token, visit
	// https://issues.teslamotors.com/secure/ViewProfile.jspa?
	//    selectedTab=com.atlassian.pats.pats-plugin:myjira-user-personal-access-tokens
	flag.StringVar(&result.Jira.Token, flagJiraToken, "",
		fmt.Sprintf("access token for the given Jira domain (overrides env var %s)", envJiraToken))

	flag.BoolVar(&result.NoTokenEcho, flagNoEcho,
		false, fmt.Sprintf("don't echo the value of tokens (over-the-shoulder security)"))

	flag.Parse()

	// All the arguments should be usernames.
	result.User = flag.Args()

	if result.Gh.Token == "" {
		result.Gh.Token = os.Getenv(envGhToken)
	}
	if result.Jira.Token == "" {
		result.Jira.Token = os.Getenv(envJiraToken)
	}

	if result.Gh.ClientId == "" {
		result.Gh.ClientId, err = determineClientIdFromDomain(result.Gh.Domain)
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
		prefix, flagGhToken, token,
		prefix, envGhToken, token,
		prefix, flagNoEcho)
}
