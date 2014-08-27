package edit

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

func (s *EditServer) NewExtract(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		languageOptions, err := args.GetLanguageOptions()
		if err != nil {
			return err
		}

		form := &server.Form{
			MustLogIn: true,
			Header:    "Create a new extract",
			Hint:      "extract_creation_hint",
			Submit:    "Save new extract",
			Fields: []*server.FormField{{
				Name:     "Slug",
				Type:     server.InputText,
				Property: "Url Slug",
				Hint:     "This is going to be the permanent address of your new extract on the web. Type something relevant to your extract, like the title. Note that special characters are not allowed here.",
			}, {
				Name:          "ExtractType",
				Type:          server.InputSelect,
				Property:      "Extract Type",
				Hint:          "Select the option which most accurately describes your new extract.",
				InputTemplate: server.InputSelect,
				Options:       server.ExtractTypeOptions,
			}, {
				Name:          "Language",
				Type:          server.InputSelect,
				Property:      "Language",
				Hint:          "Select one language (or language family). You can translate your extract into other languages later on.",
				InputTemplate: server.InputSelect,
				Options:       languageOptions,
				Link: &server.Link{
					Href: "mailto:support@polyglottis.org",
					Text: "(Not in the list? let us know by email!)",
				},
			}, {
				Name:     "Title",
				Type:     server.InputText,
				Property: "Title",
				Hint:     "Enter the title of your extract, in the language you selected above.",
			}, {
				Name:     "Summary",
				Type:     server.InputTextArea,
				Property: "Summary",
				Hint:     "Enter a short summary of your extract, in the language you selected above. This will appear in search results.",
			}, {
				Name:     "Text",
				Type:     server.InputTextArea,
				Property: "Text",
				Hint:     "IMPORTANT: Enter a line break after EACH sentence. Split paragraphs by leaving a blank line.",
			}},
		}
		form.Apply(context)
		args.Data = map[string]interface{}{
			"title": i18n.Key("New Extract"),
			"form":  form,
		}
		args.Css = "form"
		return server.FormTmpl.Execute(w, args)
	})
}
