package myjira

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/monopole/snips/internal/myhttp"
	"github.com/monopole/snips/internal/types"
)

func makeMapOfRepoToIssueList(domain string, issues []issueRecord) (map[types.RepoId][]types.MyIssue, error) {
	rawMap := make(map[types.RepoId][]*issueRecord)
	seen := make(map[string]*issueRecord)
	for i := range issues {
		issue := issues[i]
		if _, ok := seen[issue.Id]; ok {
			fmt.Fprintf(os.Stderr, "seem to be doubling up on %q\n", issue.Id)
			// ignore it - but maybe we should complain that something non-unique came back from search
			continue
		}
		seen[issue.Id] = &issue
		id := types.RepoId{
			// Name is something like 'microsoft developers'
			Org: issue.Fields.Project.Name,
			// Key is something like MSFT
			// The URL we want is https://issues.acmecorp.com/projects/MSFT/issues
			Name: issue.Fields.Project.Key,
		}
		rawMap[id] = append(rawMap[id], &issue)
	}
	var err error
	result := make(map[types.RepoId][]types.MyIssue)
	for id, jiraIssues := range rawMap {
		lst := make([]types.MyIssue, len(jiraIssues))
		for i := range jiraIssues {
			lst[i], err = convertJiraIssueToGhIssue(domain, id, jiraIssues[i])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Trouble with %+v\n", jiraIssues[i])
				return nil, err
			}
		}
		sort.Slice(lst, func(i, j int) bool {
			return lst[i].Updated.After(lst[j].Updated)
		})
		result[id] = lst
	}
	return result, nil
}

func convertJiraIssueToGhIssue(domain string, id types.RepoId, rec *issueRecord) (types.MyIssue, error) {
	updated, err := time.Parse(types.DateFormatJiraIssue, rec.Fields.Updated)
	if err != nil {
		return types.MyIssue{}, fmt.Errorf("trouble parsing 'updated' time field; %w", err)
	}
	num, err := getIssueNumber(id, rec.Key)
	if err != nil {
		return types.MyIssue{}, err
	}
	return types.MyIssue{
		RepoId:  id,
		Number:  num,
		Title:   rec.Fields.Summary,
		HtmlUrl: myhttp.Scheme + domain + "/browse/" + rec.Key,
		Updated: updated,
	}, nil
}

func getIssueNumber(id types.RepoId, raw string) (int, error) {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return 0, fmt.Errorf("expected something like PLM-1234, but have %q", raw)
	}
	if parts[0] != id.Name {
		log.Fatal(fmt.Errorf("expected %q, but got %q from issue key %s", id.Name, parts[0], raw).Error())
		// return 0, fmt.Errorf("expected %q, but got %q from issue key %s", id.Name, parts[0], raw)
	}
	return strconv.Atoi(parts[1])
}
