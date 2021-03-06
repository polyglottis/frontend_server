// Package server contains helper functions and structs for the frontend server.
package server

import (
	"fmt"

	localizer "github.com/polyglottis/frontend_server/i18n"
	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

func GetTmplArgs(context *frontend.Context) (*TmplArgs, error) {
	localizer := localizer.NewLocalizer(context)
	return &TmplArgs{
		Context:      context,
		Css:          "error",
		Localizer:    localizer,
		AngularLocal: config.Get().AngularLocal,
	}, nil
}

type TmplArgs struct {
	Data                map[string]interface{}
	Css                 string // "extract" (default), "form", or other .scss file
	Angular             bool   // angular script
	AngularLocal        bool   // angular local instead of CDN
	Context             *frontend.Context
	Description         i18n.Key
	DescriptionLitteral string // fallback if i18n.Key is not defined
	localizer.Localizer
}

func (a *TmplArgs) CanonicalUrl() string {
	return a.Context.ProtocolAndHost() + a.Context.Url
}

func (a *TmplArgs) AngularVersion() string { return "1.2.23" }

func (a *TmplArgs) Title() (string, error) {
	if title, ok := a.Data["title"].(string); ok {
		if title == "home_page" {
			return a.title(a.GetText("Polyglottis"), false), nil
		} else {
			return a.title(title, true), nil
		}
	} else if titleKey, ok := a.Data["title"].(i18n.Key); ok {
		return a.title(a.GetText(titleKey), true), nil
	} else {
		return "", fmt.Errorf("Expecting string or i18n.Key (page title) in args.Data[\"title\"].")
	}
}

func (a *TmplArgs) LoggedIn() bool {
	return a.Context.LoggedIn()
}

func (a *TmplArgs) UserName() string {
	return string(a.Context.User)
}

func (a *TmplArgs) title(t string, addBrand bool) string {
	if addBrand {
		return t + " - Polyglottis"
	} else {
		return t
	}
}

type nestedArgs struct {
	Data    interface{}
	Context *frontend.Context
	localizer.Localizer
}

func (a *TmplArgs) Nest(data interface{}) *nestedArgs {
	return &nestedArgs{
		Data:      data,
		Context:   a.Context,
		Localizer: a.Localizer,
	}
}

func (a *nestedArgs) Nest(data interface{}) *nestedArgs {
	return &nestedArgs{
		Data:      data,
		Context:   a.Context,
		Localizer: a.Localizer,
	}
}
