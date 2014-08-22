package user

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

var passwordSentTmpl = templates.Parse("user/templates/password_sent.html")

func (s *UserServer) PasswordSent(context *frontend.Context) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		args.Data = map[string]interface{}{
			"title": i18n.Key("Password Sent"),
			"email": context.Email,
		}
		args.Css = "message"
		return passwordSentTmpl.Execute(w, args)
	})
}
