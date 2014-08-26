package frontend_server

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type Server struct{}

var errTmpl = templates.Parse("templates/error.html")

func (s *Server) Error(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		args.Data = map[string]interface{}{
			"title": i18n.Key("Error"),
			"error": i18n.Key("Internal server error."),
		}
		return errTmpl.Execute(w, args)
	})
}

func (s *Server) NotFound(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		args.Data = map[string]interface{}{
			"error": i18n.Key("Page not found."),
			"title": i18n.Key("Not Found"),
		}
		return errTmpl.Execute(w, args)
	})
}
