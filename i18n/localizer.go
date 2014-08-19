package i18n

import (
	goi18n "github.com/nicksnyder/go-i18n/i18n"

	"github.com/polyglottis/platform/config"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
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
	GetText(i18n.Key) string
}

type localizer struct {
	locale language.Code
	tFunc  goi18n.TranslateFunc
}

func NewLocalizer(code language.Code) Localizer {
	locale := "en-us"
	if code != "en" {
		locale = string(code)
	}
	T, _ := goi18n.Tfunc(locale)
	return &localizer{
		locale: code,
		tFunc:  T,
	}
}

func (loc *localizer) GetText(key i18n.Key) string {
	return loc.tFunc(string(key))
}
