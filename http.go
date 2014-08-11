package frontend_server

import (
	"bytes"
	"io"
	"log"

	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

type Server struct{}

func New() *Server {
	return &Server{}
}

func getLocalizer(lang language.Code) (i18n.Localizer, error) {
	return i18n.NewLocalizer(lang), nil
}

var errTmpl = templates.Parse("templates/base", "templates/error")

type tmplArgs struct {
	Data map[string]interface{}
	i18n.Localizer
}

func (a *tmplArgs) GetKey(k string) i18n.Key {
	return a.Data[k].(i18n.Key)
}

func (s *Server) call(context *frontend.Context, f func(io.Writer, *tmplArgs) error) (answer []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from", r)
			err = r.(error)
			return
		}
	}()
	localizer, err := getLocalizer(context.Locale)
	if err != nil {
		return nil, err
	}
	args := &tmplArgs{
		Localizer: localizer,
	}

	buffer := new(bytes.Buffer)
	err = f(buffer, args)
	if err != nil {
		log.Println(err)
		buffer.Reset()
		args.Data = map[string]interface{}{
			"error": i18n.Key("Internal server error."),
		}
		err = errTmpl.Execute(buffer, args)
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

var homeTmpl = templates.Parse("templates/base", "templates/home")

func (s *Server) Home(context *frontend.Context) ([]byte, error) {
	return s.call(context, func(w io.Writer, args *tmplArgs) error {
		return homeTmpl.Execute(w, args)
	})
}

func (s *Server) NotFound(context *frontend.Context) ([]byte, error) {
	return s.call(context, func(w io.Writer, args *tmplArgs) error {
		args.Data = map[string]interface{}{
			"error": i18n.Key("Page not found."),
		}
		return errTmpl.Execute(w, args)
	})
}

var extractTmpl = templates.Parse("templates/base", "templates/extract")

func (s *Server) Extract(context *frontend.Context, extract *content.Extract) ([]byte, error) {
	return s.call(context, func(w io.Writer, args *tmplArgs) error {
		args.Data = map[string]interface{}{
			"flavorA": extract.Flavors[0],
		}
		return extractTmpl.Execute(w, args)
	})
}
