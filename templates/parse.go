package templates

import (
	"html/template"

	"github.com/polyglottis/platform/config"
)

func Parse(files ...string) *template.Template {
	// always include base template
	files = append([]string{"templates/base"}, files...)

	root := config.Get().TemplateRoot

	full := make([]string, len(files))
	for i, f := range files {
		full[i] = root + "/" + f + ".html"
	}
	return template.Must(template.ParseFiles(full...))
}
