package user

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type UserServer struct{}

var formTmpl = templates.Parse("templates/form")

func (s *UserServer) SignUp(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		form := &server.Form{
			Header: "Create your personal account",
			Submit: "Create an account",
			Fields: []*server.FormField{{
				Name:     "User",
				Type:     server.InputText,
				Property: "Username",
				Hint:     "This will be your username.",
			}, {
				Name:     "Email",
				Type:     server.InputText,
				Property: "Email Address",
				Hint:     "You will receive emails when they are big news to Polyglottis. We promise not to share your email with anyone.",
			}, {
				Name:     "Password",
				Type:     server.InputPassword,
				Property: "Password",
				Hint:     "Use at least eight characters.",
			}, {
				Name:     "PasswordConfirm",
				Type:     server.InputPassword,
				Property: "Confirm your password",
			}},
		}
		form.Apply(context)
		args.Data = map[string]interface{}{
			"title": i18n.Key("Account Creation"),
			"form":  form,
		}
		return formTmpl.Execute(w, args)
	})
}
