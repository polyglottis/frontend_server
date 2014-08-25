package user

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type UserServer struct{}

func (s *UserServer) SignUp(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		form := &server.Form{
			Header: "Create your personal account",
			Submit: "Create an account",
			Class:  "compact",
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
			}, passwordField(), passwordConfirmField()},
		}
		form.Apply(context)
		args.Data = map[string]interface{}{
			"title": i18n.Key("Account Creation"),
			"form":  form,
		}
		args.Css = "form"
		return server.FormTmpl.Execute(w, args)
	})
}

func passwordField() *server.FormField {
	return &server.FormField{
		Name:     "Password",
		Type:     server.InputPassword,
		Property: "Password",
		Hint:     "Use at least eight characters.",
	}
}

func passwordConfirmField() *server.FormField {
	return &server.FormField{
		Name:     "PasswordConfirm",
		Type:     server.InputPassword,
		Property: "Confirm your password",
	}
}
