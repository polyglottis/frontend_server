package frontend_server

import (
	"io"

	"github.com/polyglottis/frontend_server/extract"
	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type Server struct {
	*extract.Server
}

func New() *Server {
	return &Server{
		Server: extract.NewServer(),
	}
}

var homeTmpl = templates.Parse("templates/home")

func (s *Server) Home(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		args.Data = map[string]interface{}{
			"title": "home_page",
		}
		return homeTmpl.Execute(w, args)
	})
}

func (s *Server) NotFound(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		args.Data = map[string]interface{}{
			"error": i18n.Key("Page not found."),
			"title": i18n.Key("Not Found"),
		}
		return server.ErrTmpl.Execute(w, args)
	})
}
