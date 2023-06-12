package report

import (
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/monopole/snips/internal/types"
)

// mqlTimeFormat is the format used by mql's businessobject print and query commands.
const mqlTimeFormat = "1/2/2006 3:04:05 PM"

// makeFuncMap makes a string to function map for use in Go template rendering.
func makeFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"toUpper": strings.ToUpper,
		"shaSmall": func(s string) string {
			return s[0:7]
		},
		"snipDate": func(t time.Time) string {
			return t.Format(types.DayFormat2)
		},
		"prettyDateRange": func(dr *types.DayRange) string {
			return dr.PrettyRange()
		},
		"labeledIssueMap":  labeledIssueMap,
		"labeledCommitMap": labeledCommitMap,
	}
}
func labeledCommitMap(l string, m map[types.RepoId][]*types.MyCommit) interface{} {
	return &struct {
		Label string
		M     map[types.RepoId][]*types.MyCommit
	}{Label: l, M: m}
}

func labeledIssueMap(l string, m map[types.RepoId][]types.MyIssue) interface{} {
	return &struct {
		Label string
		M     map[types.RepoId][]types.MyIssue
	}{Label: l, M: m}
}

const (
	cssStyle = `
<style>
.oneIssue {
  margin-left: 10px;
}
.issueMap {
  margin-left: 20px;
  padding: 10px;
  background-color: #F8F8FF;
}
.userData {
  margin-left: 10px;
  padding-bottom: 10px;
}
</style>
`
	tmplHtmlNameIssue = "tmplHtmlNameIssue"
	tmplHtmlBodyIssue = `
{{define "` + tmplHtmlNameIssue + `" -}}
<code>{{snipDate .Updated}}</code> &nbsp; <a href="{{.HtmlUrl}}"> {{.Title}} </a>
{{- end}}
`
	tmplHtmlNameCommit = "tmplHtmlNameCommit"
	tmplHtmlBodyCommit = `
{{define "` + tmplHtmlNameCommit + `" -}}
<code>{{snipDate .Committed}}
<a href="{{.Url}}">{{shaSmall .Sha}}</a>
{{- if .Pr}} (pull/<a href="{{.Pr.HtmlUrl}}">{{.Pr.Number}}</a>){{end}}
</code>
&nbsp; {{.MessageFirstLine}}
{{- end}}
`
	tmplHtmlNameRepoToIssueMap = "tmplHtmlNameRepoToIssueMap"
	tmplHtmlBodyRepoToIssueMap = `
{{define "` + tmplHtmlNameRepoToIssueMap + `" -}}
<div class="issueMap">
{{range $repo, $list := . -}}
<h4> {{$repo}} </h4>
{{range $i, $issue := $list }}
<div class="oneIssue"> {{template "` + tmplHtmlNameIssue + `" $issue}} </div>
{{- end}}
{{- end}}
</div>
{{- end}}
`
	tmplHtmlNameRepoToCommitMap = "tmplHtmlNameRepoToCommitMap"
	tmplHtmlBodyRepoToCommitMap = `
{{define "` + tmplHtmlNameRepoToCommitMap + `" -}}
<div class="issueMap">
{{range $repo, $list := . -}}
<h4> {{$repo}} </h4>
{{range $i, $issue := $list }}
<div class="oneIssue"> {{template "` + tmplHtmlNameCommit + `" $issue}} </div>
{{- end}}
{{- end}}
</div>
{{- end}}
`

	tmplHtmlNameLabelledIssueMap = "tmplHtmlNameLabelledIssueMap"
	tmplHtmlBodyLabelledIssueMap = `
{{define "` + tmplHtmlNameLabelledIssueMap + `" -}}
{{if .M -}}
<h3> {{.Label}}: </h3>
{{template "` + tmplHtmlNameRepoToIssueMap + `" .M}}
{{- else -}}
<h3> No {{.Label}} </h3>
{{- end}}
{{- end}}
`
	tmplHtmlNameLabelledCommitMap = "tmplHtmlNameLabelledCommitMap"
	tmplHtmlBodyLabelledCommitMap = `
{{define "` + tmplHtmlNameLabelledCommitMap + `" -}}
{{if .M -}}
<h3> {{.Label}}: </h3>
{{template "` + tmplHtmlNameRepoToCommitMap + `" .M}}
{{- else -}}
<h3> No {{.Label}} </h3>
{{- end}}
{{- end}}
`
	tmplHtmlNameOrganizations = "tmplHtmlNameOrganizations"
	tmplHtmlBodyOrganizations = `
{{define "` + tmplHtmlNameOrganizations + `" -}}
<h3> Organizations </h3>
<ul>
{{range . -}} <li>{{if .Name}}{{.Name}} &nbsp; {{end}} {{.Login}}</li>
{{end}}
</ul>
{{- end}}
`
	tmplHtmlNameUser = "tmplHtmlNameUser"
	tmplHtmlBodyUser = `
{{define "` + tmplHtmlNameUser + `" -}}
<h2> {{.Name}} (<em>{{if .Email}}{{.Email}}{{else}}{{.Login}}{{end}}</em>)</h2>
<div class="userData">
{{if .Orgs}}
  {{template "` + tmplHtmlNameOrganizations + `" .Orgs}}
{{else}}
  <h3> no organizations </h3>
{{end}}
{{template "` + tmplHtmlNameLabelledIssueMap + `" (labeledIssueMap "Issues Created" .IssuesCreated)}}
{{template "` + tmplHtmlNameLabelledIssueMap + `" (labeledIssueMap "Issues Closed" .IssuesClosed)}}
{{template "` + tmplHtmlNameLabelledIssueMap + `" (labeledIssueMap "PRs Reviewed" .PrsReviewed)}}
{{template "` + tmplHtmlNameLabelledCommitMap + `" (labeledCommitMap "Commits" .Commits)}}
</div>
<hr>
{{end}}
`
	tmplHtmlNameSnipsMain = "tmplHtmlNameSnipsMain"
	tmplHtmlBodySnipsMain = `
{{define "` + tmplHtmlNameSnipsMain + `" -}}
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>` +
		cssStyle + `
  </head>
  <body>
    <h1>{{.Title}}</h1>
    <p><em> {{ prettyDateRange .Dr }} </em></p>
    {{range .Users -}}
      <div>{{ template "` + tmplHtmlNameUser + `" . -}}</div>
    {{- else -}}
      <p><strong> no users </strong></p>
    {{- end}}
  </body>
</html>
{{- end}}
`
)

// makeHtmlTemplate returns a parsed template for reporting a diff.
func makeHtmlTemplate() *template.Template {
	return template.Must(
		template.New("main").Funcs(makeFuncMap()).Parse(
			tmplHtmlBodyIssue + tmplHtmlBodyCommit + tmplHtmlBodyOrganizations +
				tmplHtmlBodyRepoToIssueMap + tmplHtmlBodyRepoToCommitMap +
				tmplHtmlBodyLabelledIssueMap + tmplHtmlBodyLabelledCommitMap +
				tmplHtmlBodyUser + tmplHtmlBodySnipsMain))
}

func WriteHtmlReport(w io.Writer, r *types.Report) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameSnipsMain, r)
}

func WriteHtmlIssue(w io.Writer, r *types.MyIssue) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameIssue, r)
}

func WriteHtmlCommit(w io.Writer, c *types.MyCommit) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameCommit, c)
}

func WriteHtmlLabelledIssueMap(w io.Writer, l string, m map[types.RepoId][]types.MyIssue) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameLabelledIssueMap, labeledIssueMap(l, m))
}

func WriteHtmlLabelledCommitMap(w io.Writer, l string, m map[types.RepoId][]*types.MyCommit) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameLabelledCommitMap, labeledCommitMap(l, m))
}

func WriteHtmlUser(w io.Writer, r *types.MyUser) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameUser, r)
}
