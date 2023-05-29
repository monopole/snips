package report

import (
	"fmt"
	"github.com/monopole/snips/internal/types"
)

func MyPrint(dr *types.DayRange, u *types.MyUser) {
	fmt.Printf("# %s %s\n", u.Name, u.Company)
	fmt.Printf("## Organizations\n")
	for i, organization := range u.Orgs {
		fmt.Printf("  %v. %v\n", i+1, organization.Login)
	}
	fmt.Printf("## Theme for %s\n", dr.PrettyRange())
	printRepoToIssueMap("Issues created", u.IssuesCreated)
	printRepoToIssueMap("Issues fixed/closed", u.IssuesClosed)
	printRepoToIssueMap("Issues commented", u.IssuesCommented)
	printRepoToIssueMap("PRs Merged", u.PrsMerged)
	printRepoToIssueMap("PRs Reviewed", u.PrsReviewed)
}

func printRepoToIssueMap(title string, m map[types.RepoName][]types.MyIssue) {
	if len(m) < 1 {
		return
	}
	fmt.Printf("\n## %s\n\n", title)
	for repo, lst := range m {
		if len(lst) < 1 {
			continue
		}
		fmt.Printf("#### %s\n\n", repo)
		for _, issue := range lst {
			printIssue(&issue)
		}
		fmt.Println()
	}
	fmt.Println()
}

func printIssue(issue *types.MyIssue) {
	fmt.Printf(
		" - %s [%s](%s)\n",
		issue.Updated.Format(types.DayFormat2),
		issue.Title,
		issue.HtmlUrl)
}
