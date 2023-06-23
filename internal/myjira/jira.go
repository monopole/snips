package myjira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/monopole/snips/internal/types"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/monopole/snips/internal/myhttp"
	"github.com/monopole/snips/internal/pgmargs"
)

const (
	// v3 is in beta, v2 is in production.  Use the latter.
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/#version
	// https://developer.atlassian.com/cloud/jira/platform/rest/v2/intro/#version
	//
	// Issue Search
	// https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-search/#api-rest-api-2-search-post
	searchEndpoint = "rest/api/2/search"

	maxResult      = 10
	maxMaxResult   = 10000
	standardFields = "id,key,summary,resolution,labels,assignee,reporter,project,description,creator,updated"
)

type jiraBoss struct {
	htCl *http.Client
	args *pgmargs.Args
}

func MakeJiraBoss(htCl *http.Client, args *pgmargs.Args) *jiraBoss {
	return &jiraBoss{
		htCl: htCl,
		args: args,
	}
}

// RANDOM jql notes
// const jql = "order by created DESC"
//
// Find all issues that were created by jregan:
// creator = jregan and created >= "2021/01/01" and created <= "2021/12/30"
//
// issuefunction in commented (" by jregan after 2023/03/01") and issuefunction in commented (" by jregan before 2023/04/01")
//
//
// the creator cannot change, but the reporter can change.  so maybe use reporter
// see :  https://support.atlassian.com/jira-software-cloud/docs/jql-fields/
//    resolved during ("date", "date")
//
//  assignee = jsmith
//   created during ("2011/02/02","2011/02/02")
//   updated during ("2011/02/02","2011/02/02")
//   reporter not in (Jack,Jill,John) and assignee not in (Jack,Jill,John)
//   duedate = empty order by created
// Comment ~ "\"text\"" and issuefunction in commented (" by user@domain.com")
// status changed by "Username" and updated> startOfDay("-1")
// status WAS "Resolved" BY jsmith DURING ("2011/02/02","2011/02/02")
//status IN ("To Do", "In Progress", "Closed")
// ...AND (issueFunction in commented("by [user]") OR issueFunction in commented("by [user]") OR issueFunction in commented("by [user]"))...

func makeIssuesCreatedJql(user string, dayRange *types.DayRange) string {
	// https://issues.teslamotors.com/browse/DESOS-625 was created at 2023-06-08 13:47
	// The following queries yield this issue
	//   created > '2023/06/07' and created < '2023/06/09'",
	//   created > '2023/06/08' and created < '2023/06/09'",
	// but this does not
	//   created >= '2023/06/08' and created <= '2023/06/08'
	return fmt.Sprintf(
		"creator = %s and created >= '%s' and created < '%s'",
		user,
		dayRange.StartAsTime().Format(types.DayFormatJira),
		// Add day really means just count up through to midnight on the last day.
		dayRange.EndAsTime().AddDate(0, 0, 1).Format(types.DayFormatJira),
	)
}

func makeJiraRequest(jql string) issueSearchRequest {
	return issueSearchRequest{
		Jql:        jql,
		MaxResults: maxResult,
		StartAt:    0,
		Fields:     strings.Split(standardFields, ","),
		Expand:     nil, // []string{"renderedFields", "names"}
	}
}

func (jb *jiraBoss) DoSearch(users []*types.MyUser) error {
	for _, u := range users {
		if err := jb.searchOneUser(u); err != nil {
			return err
		}
	}
	return nil
}

func (jb *jiraBoss) searchOneUser(user *types.MyUser) error {
	fmt.Println(user.Login)
	req := makeJiraRequest(makeIssuesCreatedJql(user.Login, jb.args.DateRange))
	loc, err := url.Parse(myhttp.Scheme + jb.args.Jira.Domain + "/" + searchEndpoint)
	if err != nil {
		return err
	}
	var issues []issueRecord
	for {
		var resp *issueSearchResponse
		resp, err = jb.doJiraRequest(loc, req)
		if err != nil {
			return err
		}
		if len(resp.Issues) == 0 {
			break
		}
		issues = append(issues, resp.Issues...)
		req.StartAt += len(resp.Issues)
		if req.StartAt > maxMaxResult {
			break
		}
	}

	newList := make([]types.MyIssue, len(issues))
	for i, jiraIssue := range issues {
		newList[i] = convertJiraIssueToGhIssue(jiraIssue)
	}
	user.IssuesCreated = make(map[types.RepoId][]types.MyIssue)
	user.IssuesCreated[jiraRepoId] = newList
	return nil
}

var jiraRepoId = types.RepoId{
	Org:  "jiraO",
	Repo: "jiraR",
}

func convertJiraIssueToGhIssue(rec issueRecord) types.MyIssue {
	return types.MyIssue{
		RepoId:  jiraRepoId,
		Number:  0,
		Title:   rec.Fields.Summary,
		HtmlUrl: "",
		Updated: time.Time{},
	}
}

func (jb *jiraBoss) doJiraRequest(loc *url.URL, req issueSearchRequest) (resp *issueSearchResponse, err error) {
	var (
		ans  io.ReadCloser
		body []byte
	)
	body, err = json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("trouble marshaling data from request; %w", err)
	}
	ans, err = jb.sendPost(loc, bytes.NewBuffer(body), jb.args.Jira.Token)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(ans)
	if err != nil {
		return nil, fmt.Errorf("ReadAll failure: %w", err)
	}
	resp = &issueSearchResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, fmt.Errorf("trouble unmarshaling data from response; %w", err)
	}
	return
}

func (jb *jiraBoss) sendPost(loc *url.URL, body io.Reader, token string) (ans io.ReadCloser, err error) {
	const debug = false
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest(http.MethodPost, loc.String(), body)
	if err != nil {
		return
	}
	req.Header.Set(myhttp.HeaderAccept, myhttp.ContentTypeJson)
	req.Header.Set(myhttp.HeaderContentType, myhttp.ContentTypeJson)
	req.Header.Set(myhttp.HeaderAAuthorization, fmt.Sprintf("Bearer: %s", token))
	resp, err = jb.htCl.Do(req)
	if err != nil {
		return
	}
	if debug {
		myhttp.PrintResponse(resp, myhttp.PrArgs{Headers: true, Body: true})
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code %d", resp.StatusCode)
		return
	}
	return resp.Body, nil
}
