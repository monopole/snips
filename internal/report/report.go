package report

import (
	"fmt"
	"github.com/monopole/snips/internal/types"
)

func PrintReport(r *types.Report) {
	fmt.Printf("# %s\n\n", r.Title)
	fmt.Printf("_%s_\n\n", r.Dr.PrettyRange())
	for i := range r.Users {
		printUser(r.Users[i])
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
	printRepoToIssueMap("PRs Reviewed", u.Login, u.PrsReviewed)
	printRepoToCommitMap("Commits", u.Login, u.Commits)
}

func printRepoToIssueMap(title string, login string, m map[types.RepoId][]types.MyIssue) {
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

func printRepoToCommitMap(title string, login string, m map[types.RepoId][]*types.MyCommit) {
	if len(m) < 1 {
		return
	}
	fmt.Printf("\n### %s (%s)\n\n", title, login)
	for repo, lst := range m {
		if len(lst) < 1 {
			continue
		}
		fmt.Printf("#### %s\n\n", repo)
		for _, c := range lst {

			fmt.Printf(
				" - `%s %s ` [%s](%s)",
				c.Committed.Format(types.DayFormat2),
				c.Sha[0:6],
				c.MessageFirstLine,
				c.Url,
			)
			if c.Pr != nil {
				fmt.Printf(
					" (pull/[%d](%s)",
					c.Pr.Number,
					c.Pr.HtmlUrl,
				)
				if c.Pr.Title != c.MessageFirstLine {
					fmt.Print(" ", c.Pr.Title)
				}
				fmt.Print(")")
			}
			fmt.Println()
		}
		fmt.Println()
	}
	fmt.Println()
}

func printIssue(issue *types.MyIssue) {
	fmt.Printf(
		" - `%s ` [%s](%s)\n",
		issue.Updated.Format(types.DayFormat2),
		issue.Title,
		issue.HtmlUrl)
}
