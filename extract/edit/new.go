package edit

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/content"
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
			Header: "Create a new extract",
			Submit: "Save new extract",
			Fields: []*server.FormField{{
				Name:     "Slug",
				Type:     server.InputText,
				Property: "Url Slug",
				Hint:     "Use at least five characters.",
			}, {
				Name:          "ExtractType",
				Type:          server.InputSelect,
				Property:      "Extract Type",
				Hint:          "Select the option which most accurately describes your new extract.",
				InputTemplate: server.InputSelect,
				Options:       extractTypeOptions,
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
				Hint:     "Enter a short summary of your extract, in the language you selected above. This will appear in search results only.",
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

var extractTypeOptions []*server.FormOption

func init() {
	extractTypeOptions = make([]*server.FormOption, len(content.AllExtractTypes)+1)
	extractTypeOptions[0] = server.PleaseSelect
	for i, eType := range content.AllExtractTypes {
		extractTypeOptions[i+1] = &server.FormOption{
			Value: string(eType),
			Key:   i18n.Key("ExtractType_" + eType),
		}
	}
}
