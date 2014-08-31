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
			Class:  "compact",
			Fields: []*server.FormField{{
				Name:     "User",
				Type:     server.InputText,
				Property: "Username or email",
			}, {
				Name:     "Password",
				Type:     server.InputPassword,
				Property: "Password",
				Link: &server.Link{
					Href: "/user/forgot_password",
					Text: "(I have forgotten my password)",
				},
			}},
		}
		form.Apply(context)
		args.Description = "Sign in to Polyglottis, the open platform for all languages."
		args.Data = map[string]interface{}{
			"title": i18n.Key("Sign In"),
			"form":  form,
		}
		args.Css = "form"
		return server.FormTmpl.Execute(w, args)
	})
}
