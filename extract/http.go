package extract

import (
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

type ExtractServer struct{}

// var languageSelectScript = templates.Script("extract/templates/language_select")
var flavorTmpl = templates.Parse("extract/templates/frame.html", "extract/templates/flavor.html", "extract/templates/language_select.js")

type TmplArgs struct {
	*server.TmplArgs
	languageA language.Code
	languageB language.Code
	ExtractId content.ExtractId
	Slug      string
}

func (a *TmplArgs) LanguageA() string {
	return a.languageString(a.languageA)
}

func (a *TmplArgs) LanguageB() string {
	return a.languageString(a.languageB)
}

func (a *TmplArgs) languageString(code language.Code) string {
	return a.GetText(i18n.Key("lang_" + string(code)))
}

func (a *TmplArgs) LinkEdit(which, what string) string {
	query := url.Values{}
	if len(a.ExtractId) == 0 {
		log.Println("Unable to generate edit link when extrat id is not set.")
		return ""
	}
	query.Set("a", string(a.languageA))
	if len(a.languageB) != 0 && a.languageB != language.Unknown.Code {
		query.Set("b", string(a.languageB))
	}
	switch which {
	case "a", "b":
		query.Set("focus", which)
	default:
		log.Println("Argument \"which\" should be either \"a\" or \"b\"")
	}
	return fmt.Sprintf("/extract/edit/%s/%s?%s", what, a.Slug, query.Encode())
}

func (a *TmplArgs) LinkRead() string {
	return fmt.Sprintf("/extract/%s/%s", a.Slug, string(a.languageA))
}

func NewTmplArgsExtract(tmplArgs *server.TmplArgs, e *content.Extract) *TmplArgs {
	args := &TmplArgs{TmplArgs: tmplArgs}
	if e != nil {
		args.ExtractId = e.Id
		args.Slug = e.UrlSlug
	}
	return args
}

func newTmplArgsTriples(tmplArgs *server.TmplArgs, e *content.Extract, a, b *frontend.FlavorTriple) *TmplArgs {
	args := NewTmplArgsExtract(tmplArgs, e)
	if a != nil {
		args.languageA = a.Language()
	}
	if b != nil {
		args.languageB = b.Language()
	}
	return args
}

func NewTmplArgs(tmplArgs *server.TmplArgs, e *content.Extract, a, b *content.Flavor) *TmplArgs {
	args := NewTmplArgsExtract(tmplArgs, e)
	if a != nil {
		args.languageA = a.Language
	}
	if b != nil {
		args.languageB = b.Language
	}
	return args
}

func (s *ExtractServer) Flavor(context *frontend.Context, extract *content.Extract, a, b *frontend.FlavorTriple) ([]byte, error) {
	return server.Call(context, func(w io.Writer, tmplArgs *server.TmplArgs) error {
		args := newTmplArgsTriples(tmplArgs, extract, a, b)
		f := shapeFlavor(extract.Shape(), a.Text)
		args.Angular = true
		args.Data = map[string]interface{}{
			"title":   getTitle(f),
			"flavorA": f,
		}
		args.Data["LanguageOptions"], args.Data["Selection"] = args.languageOptions(extract)
		return flavorTmpl.Execute(w, args)
	})
}

func getTitle(f *Flavor) interface{} {
	if f.Title.Missing {
		return i18n.Key("No title")
	} else {
		return f.Title.Content
	}
}

type Flavor struct {
	Id       content.FlavorId
	Summary  string
	Title    String
	Language language.Code
	Blocks   []*Block
}

type Block struct {
	Id      content.BlockId
	Units   []*Unit
	Missing bool
}

type Unit struct {
	Id content.UnitId
	String
}

type String struct {
	Content string
	Missing bool
}

func shapeFlavor(shape content.ExtractShape, flavor *content.Flavor) *Flavor {
	f := &Flavor{
		Id:       flavor.Id,
		Summary:  flavor.Summary,
		Title:    String{Missing: true},
		Language: flavor.Language,
	}
	if len(flavor.Blocks) != 0 {
		u := flavor.Blocks[0][0]
		if u.BlockId == 1 && u.Id == 1 {
			f.Title = String{Content: u.Content}
		}
	}
	shape.IterateFlavorBody(flavor, func(blockId content.BlockId) {
		f.Blocks = append(f.Blocks, &Block{Id: blockId})
	}, func(blockId content.BlockId, unitId content.UnitId, u *content.Unit) {
		unit := &Unit{Id: unitId}
		if u == nil {
			unit.Missing = true
		} else {
			unit.Content = u.Content
		}
		f.Blocks[len(f.Blocks)-1].Units = append(f.Blocks[len(f.Blocks)-1].Units, unit)
	}, nil)
	return f
}

type languageOption struct {
	Code  language.Code
	Label string
	Text  []*versionOption
}
type versionOption struct {
	Id      content.FlavorId
	Comment string
	Summary string
}
type selection struct {
	LanguageA int
	LanguageB int
	TextA     int
	TextB     int
}

func newSelection() *selection {
	return &selection{LanguageA: -1, LanguageB: -1}
}

func (args *TmplArgs) languageOptions(e *content.Extract) ([]*languageOption, *selection) {
	options := make([]*languageOption, 0, len(e.Flavors))
	selected := newSelection()
	for langCode, fByType := range e.Flavors {
		if flavors, ok := fByType[content.Text]; ok {
			if args.languageA == langCode {
				selected.LanguageA = len(options)
			}
			if args.languageB == langCode {
				selected.LanguageB = len(options)
			}
			options = append(options, &languageOption{
				Code:  langCode,
				Label: args.GetLanguage(langCode),
				Text:  args.versionOptions(flavors),
			})
		}
	}
	return options, selected
}

func (args *TmplArgs) versionOptions(flavors []*content.Flavor) []*versionOption {
	options := make([]*versionOption, len(flavors))
	for i, f := range flavors {
		options[i] = &versionOption{
			Id:      f.Id,
			Comment: f.LanguageComment,
			Summary: f.Summary,
		}
	}
	return options
}
