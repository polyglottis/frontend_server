package user

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

func (s *UserServer) SignIn(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		form := &server.Form{
			Header: "Sign in",
			Submit: "Sign in",
			Fields: []*server.FormField{{
				Name:     "User",
				Type:     server.InputText,
				Property: "Username or email",
			}, {
				Name:     "Password",
				Type:     server.InputPassword,
				Property: "Password",
			}},
		}
		form.Apply(context)
		args.Data = map[string]interface{}{
			"title": i18n.Key("Sign In"),
			"form":  form,
		}
		return formTmpl.Execute(w, args)
	})
}
