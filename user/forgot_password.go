package user

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

func (s *UserServer) ForgotPassword(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		form := &server.Form{
			Header: "Password forgotten",
			Submit: "Submit",
			Fields: []*server.FormField{{
				Name:     "Email",
				Type:     server.InputText,
				Property: "Email",
				Hint:     "Enter your email address.",
			}},
		}
		form.Apply(context)
		args.Data = map[string]interface{}{
			"title": i18n.Key("Forgot your password?"),
			"form":  form,
		}
		return formTmpl.Execute(w, args)
	})
}
