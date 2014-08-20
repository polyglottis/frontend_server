// Package server contains helper functions and structs for the frontend server.
package server

import (
	"fmt"
	"log"
	"net/url"

	localizer "github.com/polyglottis/frontend_server/i18n"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

func GetTmplArgs(context *frontend.Context) (*TmplArgs, error) {
	localizer := localizer.NewLocalizer(context)
	return &TmplArgs{
		Context:   context,
		Css:       "extract",
		Localizer: localizer,
	}, nil
}

type TmplArgs struct {
	Data    map[string]interface{}
	Css     string // "extract" (default), "form", or other .scss file
	Context *frontend.Context
	localizer.Localizer
}

func (a *TmplArgs) GetKey(k string) (i18n.Key, error) {
	if a.Data[k] == nil {
		return "", fmt.Errorf("Expecting i18n.Key in args.Data[\"%s\"] but found nil.", k)
	}
	if key, ok := a.Data[k].(i18n.Key); ok {
		return key, nil
	} else {
		return "", fmt.Errorf("args.Data[\"%s\"] should be of type i18n.Key")
	}
}

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
	return len(a.Context.User) != 0
}

func (a *TmplArgs) UserName() string {
	return string(a.Context.User)
}

func (a *TmplArgs) LanguageA() string {
	return a.languageString(a.Context.LanguageA)
}

func (a *TmplArgs) LanguageB() string {
	return a.languageString(a.Context.LanguageB)
}

func (a *TmplArgs) languageString(code language.Code) string {
	return a.GetText(i18n.Key("lang_" + string(code)))
}

func (a *TmplArgs) title(t string, addBrand bool) string {
	if addBrand {
		return t + " - Polyglottis"
	} else {
		return t
	}
}

func (a *TmplArgs) LinkEdit(which, what string) string {
	query := url.Values{}
	if len(a.Context.ExtractId) == 0 {
		log.Println("Unable to generate edit link when extrat id is not set.")
		return ""
	}
	query.Set("id", string(a.Context.ExtractId))
	query.Set("a", string(a.Context.LanguageA))
	if len(a.Context.LanguageB) != 0 {
		query.Set("b", string(a.Context.LanguageB))
	}
	switch which {
	case "a", "b":
		query.Set("focus", which)
	default:
		log.Println("Argument \"which\" should be either \"a\" or \"b\"")
	}
	return fmt.Sprintf("/extract/edit/%s?%s", what, query.Encode())
}

type nestedArgs struct {
	Data interface{}
	localizer.Localizer
}

func (a *TmplArgs) Nest(data interface{}) *nestedArgs {
	return &nestedArgs{
		Data:      data,
		Localizer: a.Localizer,
	}
}

func (a *nestedArgs) Nest(data interface{}) *nestedArgs {
	return &nestedArgs{
		Data:      data,
		Localizer: a.Localizer,
	}
}
