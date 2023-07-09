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

type DomainAndRepo struct {
	Dgh string
	Rid types.RepoId
}

func (dr DomainAndRepo) HRef() string {
	if strings.Contains(dr.Dgh, "github") {
		// Try to make a GitHub link.
		return dr.Dgh + "/" + dr.Rid.String()
	}
	// Try to make a Jira link.
	return dr.Dgh + "/projects/" + dr.Rid.Name + "/issues"
}

// makeFuncMap makes a string to function map for use in Go template rendering.
func makeFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"toUpper": strings.ToUpper,
		"shaSmall": func(s string) string {
			return s[0:7]
		},
		"snipDate": func(t time.Time) string {
			return t.Format(types.DayFormatHuman)
		},
		"prettyDateRange": func(dr *types.DayRange) string {
			return dr.PrettyRange()
		},
		"labeledIssueSet":  labeledIssueSet,
		"labeledCommitMap": labeledCommitMap,
		"mapTotalCommits": func(m map[types.RepoId][]*types.MyCommit) int {
			c := 0
			for _, v := range m {
				c += len(v)
			}
			return c
		},
		"bigEnough": func(s int) bool {
			return s > 5
		},
		"domainsAndUser": func(dGh string, dJira string, u *types.MyUser) interface{} {
			return &struct {
				Dgh   string
				Djira string
				U     *types.MyUser
			}{Dgh: dGh, Djira: dJira, U: u}
		},
		"domainAndCommitMap": func(dGh string, m map[types.RepoId][]*types.MyCommit) interface{} {
			return &struct {
				Dgh string
				M   map[types.RepoId][]*types.MyCommit
			}{Dgh: dGh, M: m}
		},
		"domainAndRepo": func(dGh string, rid types.RepoId) interface{} {
			return &DomainAndRepo{
				Dgh: dGh,
				Rid: rid,
			}
		},
		"domainAndOrgs": func(dGh string, o []types.MyGhOrg) interface{} {
			return &struct {
				Dgh    string
				GhOrgs []types.MyGhOrg
			}{Dgh: dGh, GhOrgs: o}
		},
	}
}

func labeledCommitMap(l string, dGh string, m map[types.RepoId][]*types.MyCommit) interface{} {
	return &struct {
		Label string
		Dgh   string
		M     map[types.RepoId][]*types.MyCommit
	}{Label: l, Dgh: dGh, M: m}
}

func labeledIssueSet(l string, iSet *types.IssueSet) interface{} {
	return &struct {
		Label string
		ISet  *types.IssueSet
	}{Label: l, ISet: iSet}
}

const (
	cssStyle = `
<style>
.oneIssue {
  margin-left: 10px;
}
.issueMap {
  margin-left: 20px;
  margin-top: 1px;
  padding-left: 10px;
  padding-top: 0px;
  padding-bottom: 2px;
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
<a href="https://{{.HRef}}"> {{.Rid}} </a>
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
	tmplHtmlNameIssueSet = "tmplHtmlNameIssueSet"
	tmplHtmlBodyIssueSet = `
{{define "` + tmplHtmlNameIssueSet + `" -}}
<div class="issueMap">
{{range $repo, $list := .Groups -}}
<h4> {{template "` + tmplHtmlNameRepoLink + `" domainAndRepo $.Domain $repo}}  {{template "` + tmplHtmlNameCount + `" len $list}}</h4>
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
<h4> {{template "` + tmplHtmlNameRepoLink + `" domainAndRepo $.Dgh $repo}} {{template "` + tmplHtmlNameCount + `" len $list}}</h4>
{{range $i, $issue := $list }}
<div class="oneIssue"> {{template "` + tmplHtmlNameCommit + `" $issue}} </div>
{{- end}}
{{- end}}
</div>
{{- end}}
`

	tmplHtmlNameLabeledIssueSet = "tmplHtmlNameLabeledIssueSet"
	tmplHtmlBodyLabeledIssueSet = `
{{define "` + tmplHtmlNameLabeledIssueSet + `" -}}
{{if (or (eq .ISet nil) .ISet.IsEmpty) -}}
<h3> No {{.Label}} </h3>
{{- else -}}
<h3> {{.Label}} {{template "` + tmplHtmlNameCount + `" .ISet.Count}}</h3>
{{template "` + tmplHtmlNameIssueSet + `" .ISet}}
{{- end}}
{{- end}}
`
	tmplHtmlNameLabeledCommitMap = "tmplHtmlNameLabeledCommitMap"
	tmplHtmlBodyLabeledCommitMap = `
{{define "` + tmplHtmlNameLabeledCommitMap + `" -}}
{{if .M -}}
<h3> {{.Label}} {{template "` + tmplHtmlNameCount + `" (mapTotalCommits .M)}} </h3>
{{template "` + tmplHtmlNameRepoToCommitMap + `" (domainAndCommitMap .Dgh .M)}}
{{- else -}}
<h3> No {{.Label}} </h3>
{{- end}}
{{- end}}
`
	tmplHtmlNameOrganizations = "tmplHtmlNameOrganizations"
	tmplHtmlBodyOrganizations = `
{{define "` + tmplHtmlNameOrganizations + `" -}}
<h3> Github Organizations </h3>
<ul>
{{range .GhOrgs }}<li>
<a href="https://{{$.Dgh}}/{{.Login}}"> {{if .Name}}{{.Name}} &nbsp; {{end}} {{.Login}} </a>
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
{{if .U.GhOrgs}}
  {{template "` + tmplHtmlNameOrganizations + `" domainAndOrgs .Dgh .U.GhOrgs}}
{{else}}
  <h3> no organizations </h3>
{{end}}
{{template "` + tmplHtmlNameLabeledIssueSet + `" (labeledIssueSet "Issues Created" .U.IssuesCreated)}}
{{template "` + tmplHtmlNameLabeledIssueSet + `" (labeledIssueSet "Issues Commented" .U.IssuesCommented)}}
{{template "` + tmplHtmlNameLabeledIssueSet + `" (labeledIssueSet "Issues Closed" .U.IssuesClosed)}}
{{template "` + tmplHtmlNameLabeledIssueSet + `" (labeledIssueSet "PRs Reviewed" .U.PrsReviewed)}}
{{template "` + tmplHtmlNameLabeledCommitMap + `" (labeledCommitMap "Commits" .Dgh .U.Commits)}}
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
    <title>{{if .Title}}{{.Title}}{{else}}Activity at {{.DomainGh}}{{end}}</title>` +
		cssStyle + `
  </head>
  <body>
    <h1>{{.Title}}</h1>
    <p><em> {{ prettyDateRange .Dr }} </em></p>
    {{range .Users -}}
      <div>{{ template "` + tmplHtmlNameUser + `" (domainsAndUser $.DomainGh $.DomainJira .) -}}</div>
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
				tmplHtmlBodyIssueSet + tmplHtmlBodyRepoToCommitMap +
				tmplHtmlBodyLabeledIssueSet + tmplHtmlBodyLabeledCommitMap +
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

func WriteHtmlLabeledIssueSet(w io.Writer, l string, is *types.IssueSet) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameLabeledIssueSet, labeledIssueSet(l, is))
}

func WriteHtmlLabeledCommitMap(w io.Writer, l string, m map[types.RepoId][]*types.MyCommit) error {
	return makeHtmlTemplate().ExecuteTemplate(w, tmplHtmlNameLabeledCommitMap, labeledCommitMap(l, "hoser.github.com", m))
}
