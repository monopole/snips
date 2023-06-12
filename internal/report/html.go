package report

import (
	"github.com/monopole/snips/internal/types"
	"html/template"
	"io"
	"strings"
	"time"
)

// mqlTimeFormat is the format used by mql's businessobject print and query commands.
const mqlTimeFormat = "1/2/2006 3:04:05 PM"

// makeFuncMap makes a string to function map for use in Go template rendering.
func makeFuncMap() map[string]interface{} {
	return map[string]interface{}{
		"toUpper": strings.ToUpper,
		"mqlTime": func(t time.Time) string {
			return t.Format(mqlTimeFormat)
		},
		"orForTypes": func(x []string) string {
			for i := range x {
				if strings.Index(x[i], " ") >= 0 {
					x[i] = "\"" + x[i] + "\""
				}
				x[i] = "type == " + x[i]
			}
			return strings.Join(x, " || ")
		},
		"prettyDateRange": func(dr *types.DayRange) string {
			return dr.PrettyRange()
		},
	}
}

const (
	tmplNameUser = "tmplNameUser"
	tmplBodyUser = `
{{define "` + tmplNameUser + `" -}}
<h2> {{.Name}} ({{ .Login}}) </h2>



{{end}}
`
	tmplNameSnipsMain = "tmplNameSnipsMain"
	tmplBodySnipsMain = `
{{define "` + tmplNameSnipsMain + `" -}}

<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
        <p><em> {{ prettyDateRange .Dr }} </em></p>
		{{range .Users -}}
           <div>{{ template "` + tmplNameUser + `" . -}}</div>
        {{- else -}}
           <div><strong> no users </strong></div>
        {{- end}}
	</body>
</html>

{{end}}
`
)

// makeTextTemplate returns a parsed template for reporting a diff.
func makeTextTemplate() *template.Template {
	return template.Must(
		template.New("main").Funcs(makeFuncMap()).Parse(
			tmplBodyUser + tmplBodySnipsMain))
}

func WriteHtml(w io.Writer, r *types.Report) {
	tmpl := makeTextTemplate()
	if err := tmpl.ExecuteTemplate(w, tmplNameSnipsMain, r); err != nil {
		panic(err)
	}
}
