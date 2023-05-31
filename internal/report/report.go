package report

import (
	"fmt"
	"github.com/monopole/snips/internal/types"
)

func PrintReport(title string, domain string, dr *types.DayRange, users []*types.MyUser) {
	fmt.Printf("# %s\n\n", title)
	fmt.Printf("_At %s, %s_\n\n", domain, dr.PrettyRange())
	for i := range users {
		printUser(users[i])
	}
}

func printUser(u *types.MyUser) {
	fmt.Printf("## %s (%s)\n\n", u.Name, u.Login)
	fmt.Printf("### Organizations\n")
	for i, organization := range u.Orgs {
		fmt.Printf("  %v. %v\n", i+1, organization.Login)
	}
	printRepoToIssueMap("Issues created", u.Login, u.IssuesCreated)
	printRepoToIssueMap("Issues fixed/closed", u.Login, u.IssuesClosed)
	printRepoToIssueMap("Issues commented", u.Login, u.IssuesCommented)
	printRepoToIssueMap("PRs Merged", u.Login, u.PrsMerged)
	printRepoToIssueMap("PRs Reviewed", u.Login, u.PrsReviewed)
	printRepoToIssueMap("Commits", u.Login, u.Commits)
}

func printRepoToIssueMap(title string, login string, m map[types.RepoName][]types.MyIssue) {
	if len(m) < 1 {
		return
	}
	fmt.Printf("\n### %s (%s)\n\n", title, login)
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
