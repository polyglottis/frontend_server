package user

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

func (s *UserServer) ResetPassword(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		form := &server.Form{
			Header: "Change your password",
			Submit: "Change password",
			Fields: []*server.FormField{passwordField, passwordConfirmField},
		}
		form.Apply(context)
		args.Data = map[string]interface{}{
			"title": i18n.Key("Reset Password"),
			"form":  form,
		}
		return formTmpl.Execute(w, args)
	})
}
