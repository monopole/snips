package myjira

import (
	"net/url"

	"github.com/monopole/snips/internal/myhttp"
	"github.com/monopole/snips/internal/types"
)

const (
	// v3 is in beta, v2 is in production.  Use the latter.
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/intro/#version
	// https://developer.atlassian.com/cloud/jira/platform/rest/v2/intro/#version
	//
	// Issue Search
	// https://developer.atlassian.com/cloud/jira/platform/rest/v2/api-group-issue-search/#api-rest-api-2-search-post
	searchEndpoint = "rest/api/2/search"

	maxResult    = 10
	maxMaxResult = 10000
)

func (jb *jiraBoss) doJiraSearch(jql string) (*types.IssueSet, error) {
	loc, err := url.Parse(myhttp.Scheme + jb.args.Domain + "/" + searchEndpoint)
	if err != nil {
		return nil, err
	}
	var issues []issueRecord
	req := makeJiraSearchRequest(jql)
	for {
		var resp *issueSearchResponse
		resp, err = jb.doJiraRequest(loc, req)
		if err != nil {
			return nil, err
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
	m, err := makeMapOfRepoToIssueList(jb.args.Domain, issues)
	if err != nil {
		return nil, err
	}
	return &types.IssueSet{
		Domain: jb.args.Domain,
		Groups: m,
	}, nil
}

func makeJiraSearchRequest(jql string) issueSearchRequest {
	return issueSearchRequest{
		Jql:        jql,
		MaxResults: maxResult,
		StartAt:    0,
		// Send Fields:nil to get all fields (but be mindful that you'll lose them when marshalling from JSON).
		Fields: []string{
			// id is a jira internal number with seven or so digits.
			"id",
			// key is something like PLM-25038, DESOS-234.
			"key",
			// summary is the issue summary, e.g. "users wants this blue thing to be red".
			"summary",
			// resolution is a struct describing the conditions of resolution.
			"resolution",
			// labels is a string array of labels.
			"labels",

			// assignee is a struct describing a user - name, email, displayName, etc.
			"assignee",
			// reporter is a struct describing a user - name, email, displayName, etc.
			"reporter",
			// creator is a struct describing a user - name, email, displayName, etc.
			"creator",

			// project is a struct with a key like "MSFT", a name like "microsoft developers", and avatar urls.
			// A project url takes the form: https://issues.acmecorp.com/projects/PLM/issues
			"project",

			// description is the long textual description of the issue.
			"description",

			// updated is the timestamp associated with the most recent update.
			"updated",
		},
		Expand: []string{"renderedFields", "names"},
	}
}
