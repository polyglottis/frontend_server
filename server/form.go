package server

import (
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

var FormTmpl = templates.Parse("templates/form.html")

type Form struct {
	Header    i18n.Key
	Error     i18n.Key
	Fields    []*FormField
	Extra     i18n.Key
	Submit    i18n.Key
	Class     string // css class
	MustLogIn bool
}

type FieldType string

const (
	InputText     FieldType = "input"
	InputTextArea FieldType = "textarea"
	InputPassword FieldType = "password"
	InputSelect   FieldType = "select"
)

type FormField struct {
	Name          string
	Type          FieldType
	Property      i18n.Key
	Value         string        // optional
	Error         i18n.Key      // optional
	Hint          i18n.Key      // optional
	Link          *Link         // optional
	InputTemplate FieldType     // optional
	Options       []*FormOption // mandatory only with InputTemplate==InputSelect
}

type Link struct {
	Href string
	Text i18n.Key
}

type FormOption struct {
	Value string
	Key   i18n.Key `json:",omitempty"` // preferred over Text
	Text  string   `json:",omitempty"` // only used if no key provided
}

var PleaseSelect = &FormOption{
	Key: i18n.Key("Please select"),
}

func (f *Form) Apply(c *frontend.Context) {
	if err, ok := c.Errors["FORM"]; ok {
		f.Error = err
	} else if f.MustLogIn && !c.LoggedIn() {
		f.Error = i18n.Key("warning_sign_in")
	}
	for _, field := range f.Fields {
		field.Value = c.Defaults.Get(field.Name)
		if err, ok := c.Errors[field.Name]; ok {
			field.Error = err
		}
	}
}
