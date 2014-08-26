package extract

import (
	"fmt"
	"io"
	"log"
	"net/url"

	localizer "github.com/polyglottis/frontend_server/i18n"
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
		shape := extract.Shape()
		var flavorA, flavorB *content.Flavor
		if a != nil {
			flavorA = a.Text
		}
		if b != nil {
			flavorB = b.Text
		}
		flavors := shapeFlavors(shape, flavorA, flavorB)
		args.Angular = true
		args.Data = map[string]interface{}{
			"title":   getTitle(flavors),
			"flavors": flavors,
			"HasA":    flavorA != nil,
			"HasB":    flavorB != nil,
		}
		args.Data["LanguageOptions"], args.Data["Selection"] = args.languageOptions(extract)
		return flavorTmpl.Execute(w, args)
	})
}

func getTitle(f *Flavor) interface{} {
	if f == nil {
		return i18n.Key("home_page")
	}
	if f.Title.MissingA {
		return i18n.Key("No title")
	} else {
		return f.Title.ContentA
	}
}

type Flavor struct {
	IdA, IdB             content.FlavorId
	Title                StringPair
	LanguageA, LanguageB language.Code
	Blocks               []*Block
}

type Block struct {
	Id    content.BlockId
	Units []*Unit
}

type Unit struct {
	Id content.UnitId
	StringPair
}

type StringPair struct {
	ContentA, ContentB string
	MissingA, MissingB bool
}

func shapeFlavors(shape content.ExtractShape, flavorA, flavorB *content.Flavor) *Flavor {
	if flavorA == nil && flavorB == nil {
		return nil
	}
	f := &Flavor{Title: StringPair{MissingA: true, MissingB: true}}
	if flavorA != nil {
		f.IdA = flavorA.Id
		f.LanguageA = flavorA.Language
		if len(flavorA.Blocks) != 0 {
			u := flavorA.Blocks[0][0]
			if u.BlockId == 1 && u.Id == 1 {
				f.Title.ContentA = u.Content
				f.Title.MissingA = false
			}
		}
	}
	if flavorB != nil {
		f.IdB = flavorB.Id
		f.LanguageB = flavorB.Language
		if len(flavorB.Blocks) != 0 {
			u := flavorB.Blocks[0][0]
			if u.BlockId == 1 && u.Id == 1 {
				f.Title.ContentB = u.Content
				f.Title.MissingB = false
			}
		}
	}
	shape.IterateFlavorBodies(flavorA, flavorB, func(blockId content.BlockId) {
		f.Blocks = append(f.Blocks, &Block{Id: blockId})
	}, func(blockId content.BlockId, unitId content.UnitId, uA, uB *content.Unit) {
		unit := &Unit{Id: unitId}
		if uA == nil {
			unit.MissingA = true
		} else {
			unit.ContentA = uA.Content
		}
		if uB == nil {
			unit.MissingB = true
		} else {
			unit.ContentB = uB.Content
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

type nestedArgs struct {
	Data    interface{}
	Context *frontend.Context
	localizer.Localizer
	HasA, HasB bool
}

func (a *TmplArgs) Nest(data interface{}) *nestedArgs {
	return &nestedArgs{
		Data:      data,
		Context:   a.Context,
		Localizer: a.Localizer,
		HasA:      a.Data["HasA"].(bool),
		HasB:      a.Data["HasB"].(bool),
	}
}

func (a *nestedArgs) Nest(data interface{}) *nestedArgs {
	return &nestedArgs{
		Data:      data,
		Context:   a.Context,
		Localizer: a.Localizer,
		HasA:      a.HasA,
		HasB:      a.HasB,
	}
}
