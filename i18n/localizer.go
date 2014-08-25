package i18n

import (
	goi18n "github.com/nicksnyder/go-i18n/i18n"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

var Supported = []string{
	"en-us",
	"de",
}

func init() {
	root := config.Get().TemplateRoot + "/i18n/"
	for _, locale := range Supported {
		goi18n.MustLoadTranslationFile(root + locale + ".all.json")
	}
}

type Localizer interface {
	GetText(key i18n.Key, optionalArgument ...interface{}) string
	Locale() string
}

type localizer struct {
	locale string
	tFunc  goi18n.TranslateFunc
}

func NewLocalizer(context *frontend.Context) Localizer {
	locale := "en-us"
	T, _ := goi18n.Tfunc(locale)
	return &localizer{
		locale: locale, // locale MUST be a valid goi18n.Tfunc locale!!!
		tFunc:  T,
	}
}

func (loc *localizer) GetText(key i18n.Key, arg ...interface{}) string {
	if len(arg) != 0 {
		if intArg, ok := arg[0].(int); ok {
			return loc.tFunc(string(key), intArg)
		} else if strArg, ok := arg[0].(string); ok {
			return loc.tFunc(string(key), map[string]interface{}{"value": strArg})
		} else if keyArg, ok := arg[0].(i18n.Key); ok {
			return loc.tFunc(string(key), map[string]interface{}{"value": loc.tFunc(string(keyArg))})
		}
	}
	return loc.tFunc(string(key))
}

func (loc *localizer) Locale() string {
	return loc.locale
}
