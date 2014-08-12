package templates

import (
	"html/template"
)

func Parse(files ...string) *template.Template {
	// always include base template
	files = append([]string{"templates/base"}, files...)

	full := make([]string, len(files))
	for i, f := range files {
		full[i] = "github.com/polyglottis/frontend_server/" + f + ".html"
	}
	return template.Must(template.ParseFiles(full...))
}
