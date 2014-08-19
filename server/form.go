package server

import (
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type Form struct {
	Header i18n.Key
	Error  i18n.Key
	Fields []*FormField
	Extra  i18n.Key
	Submit i18n.Key
}

type FormType string

const (
	InputText     FormType = "input"
	InputPassword FormType = "password"
)

type FormField struct {
	Name     string
	Type     FormType
	Property i18n.Key
	Value    string
	Error    i18n.Key
	Hint     i18n.Key
}

func (f *Form) Apply(c *frontend.Context) {
	if err, ok := c.Errors["FORM"]; ok {
		f.Error = err
	}
	for _, field := range f.Fields {
		field.Value = c.Defaults.Get(field.Name)
		if err, ok := c.Errors[field.Name]; ok {
			field.Error = err
		}
	}
}
