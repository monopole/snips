package html

const (
	tmplNameRepoLink = "tmplRepoLink"
	tmplBodyRepoLink = `
{{define "` + tmplNameRepoLink + `" -}}
<a href="https://{{.HRef}}"> {{.Rid}} </a>
{{- end}}
`
	tmplNameItemCount = "tmplItemCount"
	tmplBodyItemCount = `
{{define "` + tmplNameItemCount + `" -}}
<span class="itemCount"> ({{.C}} {{.N}}) </span>
{{- end}}
`
	//{{if bigEnough .}} <span class="itemCount"> ({{.}}) </span> {{end -}}

	tmplNameIssue = "tmplIssue"
	tmplBodyIssue = `
{{define "` + tmplNameIssue + `" -}}
<code>{{snipDate .Updated}}</code> &nbsp; <a href="{{.HtmlUrl}}"> {{.Title}} </a>
{{- end}}
`
	tmplNameCommit = "tmplCommit"
	tmplBodyCommit = `
{{define "` + tmplNameCommit + `" -}}
<code>{{snipDate .Committed}}
<a href="{{.Url}}">{{shaSmall .Sha}}</a>
{{- if .Pr}} (pull/<a href="{{.Pr.HtmlUrl}}">{{.Pr.Number}}</a>){{end}}
</code>
&nbsp; {{.MessageFirstLine}}
{{- end}}
`
	tmplNameIssueSet = "tmplIssueSet"
	tmplBodyIssueSet = `
{{define "` + tmplNameIssueSet + `" -}}
<div class="issueMap">
{{range $repo, $list := .Groups -}}
<h4> {{template "` + tmplNameRepoLink + `" domainAndRepo $.Domain $repo}} 
<span class="itemCount">({{len $list}} issues)</span>
</h4>
{{range $i, $issue := $list }}
<div class="oneIssue"> {{template "` + tmplNameIssue + `" $issue}} </div>
{{- end}}
{{- end}}
</div>
{{- end}}
`
	tmplNameRepoToCommitMap = "tmplRepoToCommitMap"
	tmplBodyRepoToCommitMap = `
{{define "` + tmplNameRepoToCommitMap + `" -}}
<div class="issueMap">
{{range $repo, $list := .M -}}
<h4> {{template "` + tmplNameRepoLink + `" domainAndRepo $.Dgh $repo}} 
<span class="itemCount">({{len $list}} commits)</span>
</h4>
{{range $i, $issue := $list }}
<div class="oneIssue"> {{template "` + tmplNameCommit + `" $issue}} </div>
{{- end}}
{{- end}}
</div>
{{- end}}
`

	tmplNameLabeledIssueSet = "tmplLabeledIssueSet"
	tmplBodyLabeledIssueSet = `
{{define "` + tmplNameLabeledIssueSet + `" -}}
{{if (or (eq .ISet nil) .ISet.IsEmpty) -}}
<h3> No {{.Label}} </h3>
{{- else -}}
<h3 id="{{lowerHyphen .Label}}"> {{.Label}}
<span class="itemCount">({{- .ISet.Count}} issues in {{.ISet.RepoCount}} repos)</span>
</h3>
{{template "` + tmplNameIssueSet + `" .ISet}}
{{- end}}
{{- end}}
`
	tmplNameLabeledCommitMap = "tmplLabeledCommitMap"
	tmplBodyLabeledCommitMap = `
{{define "` + tmplNameLabeledCommitMap + `" -}}
{{if .M -}}
<h3 id="{{lowerHyphen .Label}}"> {{.Label}} 
<span class="itemCount">({{mapTotalCommits .M}} commits to {{len .M}} repos)</span>
</h3>
{{template "` + tmplNameRepoToCommitMap + `" (domainAndCommitMap .Dgh .M)}}
{{- else -}}
<h3> No {{.Label}} </h3>
{{- end}}
{{- end}}
`
	tmplNameOrganizations = "tmplOrganizations"
	tmplBodyOrganizations = `
{{define "` + tmplNameOrganizations + `" -}}
<h3> Github Organizations </h3>
<ul>
{{range .GhOrgs }}<li>
<a href="https://{{$.Dgh}}/{{.Login}}"> {{if .Name}}{{.Name}} &nbsp; {{end}} {{.Login}} </a>
</li>
{{end}}
</ul>
{{- end}}
`
	tmplNameUserHighlights = "tmplUserHighlights"
	tmplBodyUserHighlights = `
{{define "` + tmplNameUserHighlights + `" -}}
<table>
<tr>
  <th> what </td>
  <th> items </th>
  <th> repos </th>
</tr>
{{template "` + tmplNameSummaryIssueSet + `" (labeledIssueSet "issues created" .U.IssuesCreated)}}
{{template "` + tmplNameSummaryIssueSet + `" (labeledIssueSet "issues commented" .U.IssuesCommented)}}
{{template "` + tmplNameSummaryIssueSet + `" (labeledIssueSet "issues closed" .U.IssuesClosed)}}
{{template "` + tmplNameSummaryIssueSet + `" (labeledIssueSet "PRs reviewed" .U.PrsReviewed)}}
{{template "` + tmplNameSummaryCommits + `" (labeledCommitMap "commits" .Dgh .U.Commits)}}
</table>

{{- end}}
`
	tmplNameSummaryIssueSet = "tmplSummaryIssueSet"
	tmplBodySummaryIssueSet = `
{{define "` + tmplNameSummaryIssueSet + `" -}}
{{if (or (eq .ISet nil) .ISet.IsEmpty) -}}
<tr> No {{.Label}} </tr>
{{- else -}}
<tr>
  <td> <a href="#{{lowerHyphen .Label}}">{{.Label}}</a></td>
  <td> {{.ISet.Count}} </td>
  <td> {{.ISet.RepoCount}} </td>
</tr>
{{- end}}
{{- end}}
`
	tmplNameSummaryCommits = "tmplSummaryCommits"
	tmplBodySummaryCommits = `
{{define "` + tmplNameSummaryCommits + `" -}}
{{if (or (eq .M nil) (eq (len .M) 0)) -}}
<tr> No {{.Label}} </tr>
{{- else -}}
<tr>
  <td> <a href="#{{lowerHyphen .Label}}">{{.Label}}</a></td>
  <td> {{mapTotalCommits .M}} </td>
  <td> {{len .M}} </td>
</tr>
{{- end}}
{{- end}}
`
	tmplNameUser = "tmplUser"
	tmplBodyUser = `
{{define "` + tmplNameUser + `" -}}
<h2> {{.U.Name}} (<em>{{if .U.Email}}{{.U.Email}}{{else}}{{.U.Login}}{{end}}</em>)</h2>
<div class="userData">
{{template "` + tmplNameUserHighlights + `" domainsAndUser .Dgh "jira" .U}}
{{if .U.GhOrgs}}
  {{template "` + tmplNameOrganizations + `" domainAndOrgs .Dgh .U.GhOrgs}}
{{else}}
  <h3> no organizations </h3>
{{end}}
{{template "` + tmplNameLabeledIssueSet + `" (labeledIssueSet "Issues Created" .U.IssuesCreated)}}
{{template "` + tmplNameLabeledIssueSet + `" (labeledIssueSet "Issues Commented" .U.IssuesCommented)}}
{{template "` + tmplNameLabeledIssueSet + `" (labeledIssueSet "Issues Closed" .U.IssuesClosed)}}
{{template "` + tmplNameLabeledIssueSet + `" (labeledIssueSet "PRs Reviewed" .U.PrsReviewed)}}
{{template "` + tmplNameLabeledCommitMap + `" (labeledCommitMap "Commits" .Dgh .U.Commits)}}
</div>
<hr>
{{end}}
`
	tmplNameSnipsMain = "tmplSnipsMain"
	tmplBodySnipsMain = `
{{define "` + tmplNameSnipsMain + `" -}}
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
      <div>{{ template "` + tmplNameUser + `" (domainsAndUser $.DomainGh $.DomainJira .) -}}</div>
    {{- else -}}
      <p><strong> no users </strong></p>
    {{- end}}
  </body>
</html>
{{- end}}
`

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
.itemCount {
  padding-left: 1em;
  color: gray;
  font-style: italic;
}
table td { width: 9em; border: 1px solid black; }
table td { text-align: end; padding-right: 1em; }
table th { text-align: end; padding-right: 1em; }
</style>
`
)
