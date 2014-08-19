package user

import (
	"bytes"
	"io"
	"text/template"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/user"
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

var passwordResetTmpl = template.Must(template.ParseFiles(config.Get().TemplateRoot + "/user/templates/password_reset.tmpl"))

func (s *UserServer) PasswordResetEmail(context *frontend.Context, a *user.Account, token string) ([]byte, error) {
	args, err := server.GetTmplArgs(context)
	if err != nil {
		return nil, err
	}
	args.Data = map[string]interface{}{
		"User":  a.Name,
		"Email": a.Email,
		"Token": token,
	}

	buffer := new(bytes.Buffer)
	err = passwordResetTmpl.Execute(buffer, args)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
