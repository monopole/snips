package search

import (
	"log"
	"sort"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/monopole/snips/internal/types"
)

func (se *Engine) findCommits(myUser *types.MyUser) (map[types.RepoId][]*types.MyCommit, error) {
	var lst1, lst2 []*types.MyCommit
	var err error
	if lst1, err = se.findPrsThenFindCommits(myUser); err != nil {
		return nil, err
	}
	if lst2, err = se.findAllCommits(myUser); err != nil {
		return nil, err
	}
	result := make(map[types.RepoId][]*types.MyCommit)
	seen := make(map[string]*types.MyCommit)
	for _, c := range append(lst1, lst2...) {
		if _, ok := seen[c.Sha]; !ok {
			seen[c.Sha] = c
			result[c.RepoId] = append(result[c.RepoId], c)
		}
	}
	for _, list := range result {
		sort.Slice(list, func(i, j int) bool {
			return list[i].Committed.After(list[j].Committed)
		})
	}
	return result, nil
}

const prLookupWait = 1 * time.Second

func (se *Engine) findPrsThenFindCommits(myUser *types.MyUser) (commits []*types.MyCommit, err error) {
	var lst []*github.Issue
	lst, err = se.searchIssues("merged", "author:%s", myUser.Login)
	if err != nil {
		return
	}
	var prsMerged map[types.RepoId][]types.MyIssue
	if prsMerged, err = makeMapOfRepoToIssueList(keepOnlyPrs.from(lst)); err != nil {
		return
	}

	for repo, prList := range prsMerged {
		for i := range prList {
			var commitsForPr []*types.MyCommit
			time.Sleep(prLookupWait)
			commitsForPr, err = se.getCommitsForPr(repo, &prList[i])
			if err == nil {
				commits = append(commits, commitsForPr...)
			} else {
				log.Printf("    Trouble with user %s, pr %s", myUser.Login, prList[i].HtmlUrl)
				log.Printf("    Error: %s", err.Error())
			}
		}
	}
	return
}

// getCommitsForPr finds commits by first finding a PR.
// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#list-commits-on-a-pull-request
func (se *Engine) getCommitsForPr(repo types.RepoId, prIssue *types.MyIssue) (result []*types.MyCommit, err error) {
	var (
		resp         *github.Response
		commits, lst []*github.RepositoryCommit
	)
	opts := makeListOptions()
	for {
		lst, resp, err = se.client.PullRequests.ListCommits(
			se.ctx, prIssue.RepoId.Org, prIssue.RepoId.Name, prIssue.Number, &opts)
		if err != nil {
			return nil, err
		}
		commits = append(commits, lst...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	result = make([]*types.MyCommit, len(lst))
	for i, c := range lst {
		result[i] = &types.MyCommit{
			RepoId:           prIssue.RepoId,
			Sha:              c.GetSHA(),
			Url:              c.GetHTMLURL(),
			MessageFirstLine: upToFirstLfOrEnd(c.GetCommit().GetMessage()),
			Committed:        c.GetCommit().GetCommitter().GetDate().Time,
			Author:           c.GetAuthor().GetLogin(),
			Pr:               prIssue,
		}
	}
	return result, nil
}

// findAllCommits finds all commits, including those not associated with a PR.
// Deliberately excludes merge commits, as they are usually made by GH to merge a PR that wasn't recently rebased.
func (se *Engine) findAllCommits(myUser *types.MyUser) (result []*types.MyCommit, err error) {
	var lst []*github.CommitResult
	lst, err = se.searchCommits("author-date", "merge:false author:%s", myUser.Login)
	if err != nil {
		return nil, err
	}
	if lookAtCommitterField := false; lookAtCommitterField {
		// Usually the author is the committer, so don't bother?
		var lst2 []*github.CommitResult
		lst2, err = se.searchCommits("committer-date", "merge:false committer:%s", myUser.Login)
		if err != nil {
			return nil, err
		}
		lst = append(lst, lst2...)
		// lst = removeDuplicateCommits(lst)
	}
	seen := make(map[string]*github.CommitResult)
	for _, c := range lst {
		// Discard duplicates.
		if _, ok := seen[c.GetSHA()]; ok {
			continue
		}
		seen[c.GetSHA()] = c
		var id types.RepoId
		id, err = yankRepoId(c)
		if err != nil {
			return nil, err
		}
		result = append(result, &types.MyCommit{
			RepoId:           id,
			Sha:              c.GetSHA(),
			Url:              c.GetHTMLURL(),
			MessageFirstLine: upToFirstLfOrEnd(c.GetCommit().GetMessage()),
			Committed:        c.GetCommit().GetCommitter().GetDate().Time,
			Author:           c.GetAuthor().GetLogin(),
			Pr:               nil, // We don't know if there's a PR via this code path.
		})
	}
	return
}
