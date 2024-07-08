package html

import (
	"github.com/monopole/snips/internal/report/common"
	"github.com/monopole/snips/internal/types"
	"html/template"
	"io"
)

func makeHtmlTemplate() *template.Template {
	return template.Must(
		template.New("main").Funcs(common.MakeFuncMap()).Parse(
			tmplBodyRepoLink +
				tmplBodyItemCount +
				tmplBodyIssue +
				tmplBodyCommit +
				tmplBodyOrganizations +
				tmplBodyIssueSet +
				tmplBodyRepoToCommitMap +
				tmplBodyLabeledIssueSet +
				tmplBodyLabeledCommitMap +
				tmplBodyUser +
				tmplBodyUserHighlights +
				tmplBodySummaryIssueSet +
				tmplBodySummaryCommits +
				tmplBodySnipsMain))
}

func WriteHtmlReport(w io.Writer, r *types.Report) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplNameSnipsMain, r)
}

func WriteHtmlIssue(w io.Writer, r *types.MyIssue) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplNameIssue, r)
}

func WriteHtmlCommit(w io.Writer, c *types.MyCommit) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplNameCommit, c)
}

func WriteHtmlLabeledIssueSet(w io.Writer, l string, is *types.IssueSet) error {
	return makeHtmlTemplate().ExecuteTemplate(
		w, tmplNameLabeledIssueSet, common.LabeledIssueSet(l, is))
}

func WriteHtmlLabeledCommitMap(
	w io.Writer, l string, m map[types.RepoId][]*types.MyCommit) error {
	return makeHtmlTemplate().ExecuteTemplate(
		w, tmplNameLabeledCommitMap,
		common.LabeledCommitMap(l, "hoser.github.com", m))
}
