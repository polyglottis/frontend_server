package extract

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type ExtractServer struct{}

var flavorTmpl = templates.Parse("extract/templates/frame.html", "extract/templates/flavor.html", "extract/templates/language_select.js")

func (s *ExtractServer) Flavor(context *frontend.Context, extract *content.Extract, a, b *frontend.FlavorTriple) ([]byte, error) {
	return server.Call(context, func(w io.Writer, tmplArgs *server.TmplArgs) error {
		args := newTmplArgsTriples(tmplArgs, extract, a, b)
		shape := extract.Shape()
		var flavorA, flavorB *content.Flavor
		if a != nil {
			flavorA = a.Text
		}
		if b != nil {
			flavorB = b.Text
		}
		flavors := shapeFlavors(shape, flavorA, flavorB)
		args.Angular = true
		args.Data = map[string]interface{}{
			"title":   getTitle(flavors),
			"flavors": flavors,
			"HasA":    flavorA != nil,
			"HasB":    flavorB != nil,
		}
		args.Data["LanguageOptions"], args.Data["Selection"] = args.languageOptions(extract)
		return flavorTmpl.Execute(w, args)
	})
}

func getTitle(f *Flavor) interface{} {
	if f == nil {
		return "home_page"
	}
	if f.Title.MissingA {
		return i18n.Key("No title")
	} else {
		return f.Title.ContentA
	}
}

