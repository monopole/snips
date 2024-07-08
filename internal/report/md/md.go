package md

import (
	"github.com/monopole/snips/internal/report/common"
	"github.com/monopole/snips/internal/types"
	"html/template"
	"io"
)

const (
	tmplNameIssue = "tmplNameIssue"
	tmplBodyIssue = `
{{define "` + tmplNameIssue + `" -}}
` + "`{{snipDate .Updated}}`" + ` [{{.Title}}]({{.HtmlUrl}})
{{- end}}
`
	tmplNameCommit = "tmplNameCommit"
	tmplBodyCommit = `
{{define "` + tmplNameCommit + `" -}}
` + "`{{snipDate .Committed}}`" + " [`{{shaSmall .Sha}}`]({{.Url}})" + `
{{- if .Pr}} (pull/[{{.Pr.Number}}]({{.Pr.HtmlUrl}})){{end}} {{.MessageFirstLine}}
{{- end}}
`
	tmplNameRepoToIssueSet = "tmplNameRepoToIssueSet"
	tmplBodyRepoToIssueSet = `
{{define "` + tmplNameRepoToIssueSet + `" -}}
{{range $repo, $list := .Groups }}
#### {{$repo}}
{{range $i, $issue := $list }}
  - {{template "` + tmplNameIssue + `" $issue}}
{{- end}}
{{end}}
{{- end}}
`
	tmplNameRepoToCommitMap = "tmplNameRepoToCommitMap"
	tmplBodyRepoToCommitMap = `
{{define "` + tmplNameRepoToCommitMap + `" -}}
{{range $repo, $list := . }}
#### {{$repo}}
{{range $i, $issue := $list }}
 - {{template "` + tmplNameCommit + `" $issue}}
{{- end}}
{{end}}
{{- end}}
`

	tmplNameLabelledIssueSet = "tmplNameLabelledIssueSet"
	tmplBodyLabelledIssueSet = `
{{define "` + tmplNameLabelledIssueSet + `" -}}
{{if .ISet -}}
### {{.Label}}:
{{template "` + tmplNameRepoToIssueSet + `" .ISet}}
{{- else -}}
### No {{.Label}}
{{- end}}
{{- end}}
`
	tmplNameLabelledCommitMap = "tmplNameLabelledCommitMap"
	tmplBodyLabelledCommitMap = `
{{define "` + tmplNameLabelledCommitMap + `" -}}
{{if .M -}}
### {{.Label}}
{{template "` + tmplNameRepoToCommitMap + `" .M}}
{{- else -}}
### No {{.Label}}
{{- end}}
{{- end}}
`
	tmplNameOrganizations = "tmplNameOrganizations"
	tmplBodyOrganizations = `
{{define "` + tmplNameOrganizations + `" -}}
### Organizations
{{range . -}} * {{if .Name}}{{.Name}} {{end}} {{.Login}}
{{end}}
{{- end}}
`
	tmplNameUser = "tmplNameUser"
	tmplBodyUser = `
{{define "` + tmplNameUser + `"}}
## {{.Name}} (_{{if .Email}}{{.Email}}{{else}}{{.Login}}{{end}}_)
{{if .GhOrgs}}
{{template "` + tmplNameOrganizations + `" .GhOrgs}}
{{else}}
### no organizations
{{end}}
{{template "` + tmplNameLabelledIssueSet + `" (labeledIssueSet "Issues Created" .IssuesCreated)}}
{{template "` + tmplNameLabelledIssueSet + `" (labeledIssueSet "Issues Closed" .IssuesClosed)}}
{{template "` + tmplNameLabelledIssueSet + `" (labeledIssueSet "PRs Reviewed" .PrsReviewed)}}
{{template "` + tmplNameLabelledCommitMap + `" (labeledCommitMap "Commits" "foo.com" .Commits)}}
---
{{end}}
`
	tmplNameSnipsMain = "tmplNameSnipsMain"
	tmplBodySnipsMain = `
{{define "` + tmplNameSnipsMain + `" -}}
# {{.Title}}
_{{ prettyDateRange .Dr }}_
{{range .Users -}}
   {{ template "` + tmplNameUser + `" . -}}
{{- else -}}
__no users__
{{- end}}
{{- end}}
`
)

func makeMdTemplate() *template.Template {
	return template.Must(
		template.New("main").Funcs(common.MakeFuncMap()).Parse(
			tmplBodyIssue + tmplBodyCommit + tmplBodyOrganizations +
				tmplBodyRepoToIssueSet + tmplBodyRepoToCommitMap +
				tmplBodyLabelledIssueSet + tmplBodyLabelledCommitMap +
				tmplBodyUser + tmplBodySnipsMain))
}

func WriteMdReport(w io.Writer, r *types.Report) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplNameSnipsMain, r)
}

func WriteMdIssue(w io.Writer, r *types.MyIssue) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplNameIssue, r)
}

func WriteMdCommit(w io.Writer, c *types.MyCommit) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplNameCommit, c)
}

func WriteMdLabelledIssueSet(w io.Writer, l string, is *types.IssueSet) error {
	return makeMdTemplate().ExecuteTemplate(
		w, tmplNameLabelledIssueSet, common.LabeledIssueSet(l, is))
}

func WriteMdLabelledCommitMap(
	w io.Writer, l string, m map[types.RepoId][]*types.MyCommit) error {
	return makeMdTemplate().ExecuteTemplate(
		w, tmplNameLabelledCommitMap, common.LabeledCommitMap(l, "d.com", m))
}
