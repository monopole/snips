package report

import (
	"github.com/monopole/snips/internal/types"
	"html/template"
	"io"
)

const (
	tmplMdNameIssue = "tmplMdNameIssue"
	tmplMdBodyIssue = `
{{define "` + tmplMdNameIssue + `" -}}
` + "`{{snipDate .Updated}}`" + ` [{{.Title}}]({{.HtmlUrl}})
{{- end}}
`
	tmplMdNameCommit = "tmplMdNameCommit"
	tmplMdBodyCommit = `
{{define "` + tmplMdNameCommit + `" -}}
` + "`{{snipDate .Committed}}`" + " [`{{shaSmall .Sha}}`]({{.Url}})" + `
{{- if .Pr}} (pull/[{{.Pr.Number}}]({{.Pr.HtmlUrl}})){{end}} {{.MessageFirstLine}}
{{- end}}
`
	tmplMdNameRepoToIssueSet = "tmplMdNameRepoToIssueSet"
	tmplMdBodyRepoToIssueSet = `
{{define "` + tmplMdNameRepoToIssueSet + `" -}}
{{range $repo, $list := .Groups }}
#### {{$repo}}
{{range $i, $issue := $list }}
  - {{template "` + tmplMdNameIssue + `" $issue}}
{{- end}}
{{end}}
{{- end}}
`
	tmplMdNameRepoToCommitMap = "tmplMdNameRepoToCommitMap"
	tmplMdBodyRepoToCommitMap = `
{{define "` + tmplMdNameRepoToCommitMap + `" -}}
{{range $repo, $list := . }}
#### {{$repo}}
{{range $i, $issue := $list }}
 - {{template "` + tmplMdNameCommit + `" $issue}}
{{- end}}
{{end}}
{{- end}}
`

	tmplMdNameLabelledIssueSet = "tmplMdNameLabelledIssueSet"
	tmplMdBodyLabelledIssueSet = `
{{define "` + tmplMdNameLabelledIssueSet + `" -}}
{{if .ISet -}}
### {{.Label}}:
{{template "` + tmplMdNameRepoToIssueSet + `" .ISet}}
{{- else -}}
### No {{.Label}}
{{- end}}
{{- end}}
`
	tmplMdNameLabelledCommitMap = "tmplMdNameLabelledCommitMap"
	tmplMdBodyLabelledCommitMap = `
{{define "` + tmplMdNameLabelledCommitMap + `" -}}
{{if .M -}}
### {{.Label}}
{{template "` + tmplMdNameRepoToCommitMap + `" .M}}
{{- else -}}
### No {{.Label}}
{{- end}}
{{- end}}
`
	tmplMdNameOrganizations = "tmplMdNameOrganizations"
	tmplMdBodyOrganizations = `
{{define "` + tmplMdNameOrganizations + `" -}}
### Organizations
{{range . -}} * {{if .Name}}{{.Name}} {{end}} {{.Login}}
{{end}}
{{- end}}
`
	tmplMdNameUser = "tmplMdNameUser"
	tmplMdBodyUser = `
{{define "` + tmplMdNameUser + `"}}
## {{.Name}} (_{{if .Email}}{{.Email}}{{else}}{{.Login}}{{end}}_)
{{if .GhOrgs}}
{{template "` + tmplMdNameOrganizations + `" .GhOrgs}}
{{else}}
### no organizations
{{end}}
{{template "` + tmplMdNameLabelledIssueSet + `" (labeledIssueSet "Issues Created" .IssuesCreated)}}
{{template "` + tmplMdNameLabelledIssueSet + `" (labeledIssueSet "Issues Closed" .IssuesClosed)}}
{{template "` + tmplMdNameLabelledIssueSet + `" (labeledIssueSet "PRs Reviewed" .PrsReviewed)}}
{{template "` + tmplMdNameLabelledCommitMap + `" (labeledCommitMap "Commits" "foo.com" .Commits)}}
---
{{end}}
`
	tmplMdNameSnipsMain = "tmplMdNameSnipsMain"
	tmplMdBodySnipsMain = `
{{define "` + tmplMdNameSnipsMain + `" -}}
# {{.Title}}
_{{ prettyDateRange .Dr }}_
{{range .Users -}}
   {{ template "` + tmplMdNameUser + `" . -}}
{{- else -}}
__no users__
{{- end}}
{{- end}}
`
)

func makeMdTemplate() *template.Template {
	return template.Must(
		template.New("main").Funcs(makeFuncMap()).Parse(
			tmplMdBodyIssue + tmplMdBodyCommit + tmplMdBodyOrganizations +
				tmplMdBodyRepoToIssueSet + tmplMdBodyRepoToCommitMap +
				tmplMdBodyLabelledIssueSet + tmplMdBodyLabelledCommitMap +
				tmplMdBodyUser + tmplMdBodySnipsMain))
}

func WriteMdReport(w io.Writer, r *types.Report) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplMdNameSnipsMain, r)
}

func WriteMdIssue(w io.Writer, r *types.MyIssue) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplMdNameIssue, r)
}

func WriteMdCommit(w io.Writer, c *types.MyCommit) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplMdNameCommit, c)
}

func WriteMdLabelledIssueSet(w io.Writer, l string, is *types.IssueSet) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplMdNameLabelledIssueSet, labeledIssueSet(l, is))
}

func WriteMdLabelledCommitMap(w io.Writer, l string, m map[types.RepoId][]*types.MyCommit) error {
	return makeMdTemplate().ExecuteTemplate(w, tmplMdNameLabelledCommitMap, labeledCommitMap(l, "d.com", m))
}
