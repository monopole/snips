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
		"mapTotalCommits": func(m map[types.RepoId][]*types.MyCommit) int {
			c := 0
			for _, v := range m {
				c += len(v)
			}
			return c
		},
		"mapTotalIssues": func(m map[types.RepoId][]types.MyIssue) int {
			c := 0
			for _, v := range m {
				c += len(v)
			}
			return c
		},
		"bigEnough": func(s int) bool {
			return s > 5
		},
		"domainAndUser": func(d string, u *types.MyUser) interface{} {
			return &struct {
				D string
				U *types.MyUser
			}{D: d, U: u}
		},
		"domainAndIssueMap": func(d string, m map[types.RepoId][]types.MyIssue) interface{} {
			return &struct {
				D string
				M map[types.RepoId][]types.MyIssue
			}{D: d, M: m}
		},
		"domainAndCommitMap": func(d string, m map[types.RepoId][]*types.MyCommit) interface{} {
			return &struct {
				D string
				M map[types.RepoId][]*types.MyCommit
			}{D: d, M: m}
		},
		"domainAndRepo": func(d string, rid types.RepoId) interface{} {
			return &struct {
				D   string
				Rid types.RepoId
			}{D: d, Rid: rid}
		},
		"domainAndOrgs": func(d string, o []types.MyOrg) interface{} {
			return &struct {
				D    string
				Orgs []types.MyOrg
			}{D: d, Orgs: o}
		},
	}
}

func labeledCommitMap(l string, d string, m map[types.RepoId][]*types.MyCommit) interface{} {
	return &struct {
		Label string
		D     string
		M     map[types.RepoId][]*types.MyCommit
	}{Label: l, D: d, M: m}
}

func labeledIssueMap(l string, d string, m map[types.RepoId][]types.MyIssue) interface{} {
	return &struct {
		Label string
		D     string
		M     map[types.RepoId][]types.MyIssue
	}{Label: l, D: d, M: m}
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
.count {
  color: gray;
  font-style: italic;
}
</style>
`

	tmplHtmlNameRepoLink = "tmplHtmlRepoLink"
	tmplHtmlBodyRepoLink = `
{{define "` + tmplHtmlNameRepoLink + `" -}}
<a href="https:{{.D}}/{{.Rid}}"> {{.Rid}} </a>
{{- end}}
`
	tmplHtmlNameCount = "tmplHtmlNameCount"
	tmplHtmlBodyCount = `
{{define "` + tmplHtmlNameCount + `" -}}
{{if bigEnough .}} <span class="count"> ({{.}}) </span> {{end -}}
{{- end}}
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
{{range $repo, $list := .M -}}
<h4> {{template "` + tmplHtmlNameRepoLink + `" domainAndRepo $.D $repo}}  {{template "` + tmplHtmlNameCount + `" len $list}}</h4>
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
{{range $repo, $list := .M -}}
<h4> {{template "` + tmplHtmlNameRepoLink + `" domainAndRepo $.D $repo}} {{template "` + tmplHtmlNameCount + `" len $list}}</h4>
{{range $i, $issue := $list }}
<div class="oneIssue"> {{template "` + tmplHtmlNameCommit + `" $issue}} </div>
{{- end}}
{{- end}}
</div>
{{- end}}
`

	tmplHtmlNameLabeledIssueMap = "tmplHtmlNameLabeledIssueMap"
	tmplHtmlBodyLabeledIssueMap = `
{{define "` + tmplHtmlNameLabeledIssueMap + `" -}}
{{if .M -}}
<h3> {{.Label}} {{template "` + tmplHtmlNameCount + `" (mapTotalIssues .M)}}: </h3>
{{template "` + tmplHtmlNameRepoToIssueMap + `" (domainAndIssueMap .D .M)}}
{{- else -}}
<h3> No {{.Label}} </h3>
{{- end}}
{{- end}}
`
	tmplHtmlNameLabeledCommitMap = "tmplHtmlNameLabeledCommitMap"
	tmplHtmlBodyLabeledCommitMap = `
{{define "` + tmplHtmlNameLabeledCommitMap + `" -}}
{{if .M -}}
<h3> {{.Label}} {{template "` + tmplHtmlNameCount + `" (mapTotalCommits .M)}}: </h3>
{{template "` + tmplHtmlNameRepoToCommitMap + `" (domainAndCommitMap .D .M)}}
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
{{range .Orgs }}<li>
<a href="https://{{$.D}}/{{.Login}}"> {{if .Name}}{{.Name}} &nbsp; {{end}} {{.Login}} </a>
</li>
{{end}}
</ul>
{{- end}}
`
	tmplHtmlNameUser = "tmplHtmlNameUser"
	tmplHtmlBodyUser = `
{{define "` + tmplHtmlNameUser + `" -}}
<h2> {{.U.Name}} (<em>{{if .U.Email}}{{.U.Email}}{{else}}{{.U.Login}}{{end}}</em>)</h2>
<div class="userData">
{{if .U.Orgs}}
  {{template "` + tmplHtmlNameOrganizations + `" domainAndOrgs .D .U.Orgs}}
{{else}}
  <h3> no organizations </h3>
{{end}}
{{template "` + tmplHtmlNameLabeledIssueMap + `" (labeledIssueMap "Issues Created" .D .U.IssuesCreated)}}
{{template "` + tmplHtmlNameLabeledIssueMap + `" (labeledIssueMap "Issues Closed" .D .U.IssuesClosed)}}
{{template "` + tmplHtmlNameLabeledIssueMap + `" (labeledIssueMap "PRs Reviewed" .D .U.PrsReviewed)}}
{{template "` + tmplHtmlNameLabeledCommitMap + `" (labeledCommitMap "Commits" .D .U.Commits)}}
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
    <title>{{if .Title}}{{.Title}}{{else}}Activity at {{.Domain}}{{end}}</title>` +
		cssStyle + `
  </head>
  <body>
    <h1>{{.Title}}</h1>
    <p><em> {{ prettyDateRange .Dr }} </em></p>
    {{range .Users -}}
      <div>{{ template "` + tmplHtmlNameUser + `" (domainAndUser $.Domain .) -}}</div>
    {{- else -}}
      <p><strong> no users </strong></p>
    {{- end}}
  </body>
</html>
{{- end}}
`
)

func makeHtmlTemplate() *template.Template {
	return template.Must(
		template.New("main").Funcs(makeFuncMap()).Parse(
			tmplHtmlBodyRepoLink + tmplHtmlBodyCount +
				tmplHtmlBodyIssue + tmplHtmlBodyCommit + tmplHtmlBodyOrganizations +
				tmplHtmlBodyRepoToIssueMap + tmplHtmlBodyRepoToCommitMap +
				tmplHtmlBodyLabeledIssueMap + tmplHtmlBodyLabeledCommitMap +
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

func WriteHtmlLabeledIssueMap(w io.Writer, l string, m map[types.RepoId][]types.MyIssue) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameLabeledIssueMap, labeledIssueMap(l, "github.com", m))
}

func WriteHtmlLabeledCommitMap(w io.Writer, l string, m map[types.RepoId][]*types.MyCommit) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameLabeledCommitMap, labeledCommitMap(l, "hoser.com", m))
}
