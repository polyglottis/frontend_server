package frontend_server

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/frontend"
)

var homeTmpl = templates.Parse("templates/home.html", "templates/home.js")

func (s *Server) Home(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, serverArgs *server.TmplArgs) error {
		args := &homeArgs{serverArgs}
		args.Data = map[string]interface{}{
			"title": "home_page",
		}
		args.Angular = true
		args.Css = "home"
		return homeTmpl.Execute(w, args)
	})
}

type homeArgs struct {
	*server.TmplArgs
}

func (a *homeArgs) ExtractTypes() []*server.FormOption {
	orig := server.ExtractTypeOptions[1:]
	translated := make([]*server.FormOption, len(orig))
	for i, o := range orig {
		translated[i] = &server.FormOption{
			Value: o.Value,
			Text:  a.GetText(o.Key),
		}
	}
	return translated
}
