package edit

import (
	"io"

	"github.com/polyglottis/frontend_server/extract"
	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

var newFlavorTmpl = templates.Parse("templates/form.html", "extract/edit/templates/new_flavor.html", "extract/edit/templates/new_flavor.js")

func (s *EditServer) NewFlavor(context *frontend.Context, e *content.Extract) ([]byte, error) {
	return server.Call(context, func(w io.Writer, serverArgs *server.TmplArgs) error {
		args := extract.NewTmplArgsExtract(serverArgs, e)
		languageOptions, err := args.GetLanguageOptions()
		if err != nil {
			return err
		}

		form := &server.Form{
			MustLogIn: true,
			Header:    "Create a new language flavor",
			Submit:    "Save",
		}
		if queryLangA := context.Query.Get("a"); len(queryLangA) != 0 {
			context.Defaults.Set("Language", queryLangA)
		}
		form.Apply(context)
		args.Angular = true
		args.Data = map[string]interface{}{
			"title":           i18n.Key("New Flavor"),
			"form":            form,
			"Language":        context.Defaults.Get("Language"),
			"LanguageOptions": languageOptions[1:],
			"Flavors":         e.Flavors,
			"errors":          context.Errors,
		}
		args.Css = "form"
		return newFlavorTmpl.Execute(w, args)
	})
}
