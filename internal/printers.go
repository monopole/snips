package internal

import (
	"fmt"

	"github.com/google/go-github/v33/github"
)

func PrintIssues(title string, issues []*github.Issue) {
	fmt.Printf("\n## %s\n\n", title)
	for repo, issueList := range MapRepoToIssueList(issues) {
		fmt.Printf("#### %s\n\n", repo)
		for _, issue := range issueList {
			printIssueLink(issue)
		}
		fmt.Println()
	}
	fmt.Println()
}

func PrintPrLink(pr *github.PullRequest) {
	fmt.Printf(" - %s [%s](%s)\n",
		pr.GetMergedAt().Format("2006-01-02"), pr.GetTitle(), pr.GetHTMLURL())
}

func printIssueLink(issue *github.Issue) {
	fmt.Printf(
		" - %s [%s](%s)\n",
		issue.GetUpdatedAt().Format("2006-01-02"), *issue.Title, *issue.HTMLURL)
}
