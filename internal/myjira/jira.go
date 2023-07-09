package myjira

import (
	"fmt"
	"net/http"
	"time"

	"github.com/monopole/snips/internal/pgmargs"
	"github.com/monopole/snips/internal/types"
)

type jiraBoss struct {
	htCl     *http.Client
	args     *pgmargs.ServiceArgs
	dayRange *types.DayRange
}

func MakeJiraBoss(htCl *http.Client, args *pgmargs.ServiceArgs, dayRange *types.DayRange) *jiraBoss {
	return &jiraBoss{
		htCl:     htCl,
		args:     args,
		dayRange: dayRange,
	}
}

func (jb *jiraBoss) DoSearch(users []*types.MyUser) (err error) {
	for _, u := range users {
		u.IssuesCreated, err = jb.doJiraSearch(makeIssuesCreatedJql(u.Login, jb.dayRange))
		if err != nil {
			return err
		}
		u.IssuesClosed, err = jb.doJiraSearch(makeIssuesClosedJql(u.Login, jb.dayRange))
		if err != nil {
			return err
		}
		u.IssuesCommented, err = jb.doJiraSearch(makeIssuesCommentedJql(u.Login, jb.dayRange))
		if err != nil {
			return err
		}
	}
	return nil
}

func makeIssuesCreatedJql(user string, dayRange *types.DayRange) string {
	// the creator cannot change, but the reporter can change.  so maybe use reporter
	// see :  https://support.atlassian.com/jira-software-cloud/docs/jql-fields/
	return fmt.Sprintf(
		"creator = %s and created >= '%s' and created < '%s'",
		user,
		dayRange.StartAsTime().Format(types.DayFormatJira),
		adjustEndDate(dayRange.EndAsTime()).Format(types.DayFormatJira),
	)
}

func makeIssuesCommentedJql(user string, dayRange *types.DayRange) string {
	return fmt.Sprintf(
		"creator != %s and issuefunction in commented (' by %s after %s') and issuefunction in commented ('by %s before %s')",
		user,
		user,
		dayRange.StartAsTime().Format(types.DayFormatJira),
		user,
		adjustEndDate(dayRange.EndAsTime()).Format(types.DayFormatJira),
	)
}

func makeIssuesClosedJql(user string, dayRange *types.DayRange) string {
	return fmt.Sprintf(
		"status WAS 'Resolved' BY %s DURING ('%s','%s')",
		user,
		dayRange.StartAsTime().Format(types.DayFormatJira),
		adjustEndDate(dayRange.EndAsTime()).Format(types.DayFormatJira),
	)
}

// adjustEndDate adds one day to the end-day so that the query range counts up through midnight on the end-day.
//
// E.g. we know that https://issues.acmecorp.com/browse/MSFT-001 was created at 2023-06-08T13:47
//
// The following queries yield this issue
//
//	created > '2023/06/07' and created < '2023/06/09'",
//	created > '2023/06/08' and created < '2023/06/09'",
//
// but this query fails:
//
//	created >= '2023/06/08' and created <= '2023/06/08'
func adjustEndDate(ed time.Time) time.Time {
	return ed.AddDate(0, 0, 1)
}
