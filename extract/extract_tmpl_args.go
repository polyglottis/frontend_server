package extract

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	localizer "github.com/polyglottis/frontend_server/i18n"
	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

type TmplArgs struct {
	*server.TmplArgs
	languageA, languageB language.Code
	textIdA, textIdB     content.FlavorId
	ExtractId            content.ExtractId
	Slug                 string
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
	if len(a.Slug) == 0 {
		log.Println("Unable to generate edit link when extrat slug is not set.")
		return ""
	}
	query := a.query(true)
	switch which {
	case "a", "b":
		query.Set("focus", which)
	default:
		log.Println("Argument \"which\" should be either \"a\" or \"b\"")
	}
	return fmt.Sprintf("/extract/edit/%s/%s?%s", what, a.Slug, query.Encode())
}

func (a *TmplArgs) query(includeA bool) url.Values {
	query := url.Values{}
	if len(a.languageA) != 0 && a.languageA != language.Unknown.Code {
		if includeA {
			query.Set("a", string(a.languageA))
		}
		if a.textIdA != 0 {
			query.Set("at", strconv.Itoa(int(a.textIdA)))
		}
	}
	if len(a.languageB) != 0 && a.languageB != language.Unknown.Code {
		query.Set("b", string(a.languageB))
		if a.textIdB != 0 {
			query.Set("bt", strconv.Itoa(int(a.textIdB)))
		}
	}
	return query
}

func (a *TmplArgs) LinkRead() string {
	query := a.query(false)
	path := fmt.Sprintf("/extract/%s/%s", a.Slug, string(a.languageA))
	if len(query) == 0 {
		return path
	}
	return path + "?" + query.Encode()
}

func (a *TmplArgs) LinkNewFlavor() string {
	query := a.query(true)
	path := fmt.Sprintf("/extract/edit/new_flavor/%s", a.Slug)
	if len(query) == 0 {
		return path
	}
	return path + "?" + query.Encode()
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
		if a.Text != nil {
			args.textIdA = a.Text.Id
		}
	}
	if b != nil {
		args.languageB = b.Language()
		if b.Text != nil {
			args.textIdB = b.Text.Id
		}
	}
	return args
}

func NewTmplArgs(tmplArgs *server.TmplArgs, e *content.Extract, a, b *content.Flavor) *TmplArgs {
	args := NewTmplArgsExtract(tmplArgs, e)
	if a != nil {
		args.languageA = a.Language
		args.textIdA = a.Id
	}
	if b != nil {
		args.languageB = b.Language
		args.textIdB = b.Id
	}
	return args
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

func (args *TmplArgs) LanguageOptions(e *content.Extract) ([]*languageOption, *selection) {
	options := make([]*languageOption, 0, len(e.Flavors))
	selected := newSelection()
	for langCode, fByType := range e.Flavors {
		if flavors, ok := fByType[content.Text]; ok {
			if args.languageA == langCode {
				selected.LanguageA = len(options)
				selected.TextA = args.indexOfSelection(flavors, args.textIdA)
			}
			if args.languageB == langCode {
				selected.LanguageB = len(options)
				selected.TextB = args.indexOfSelection(flavors, args.textIdB)
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

func (args *TmplArgs) indexOfSelection(flavors []*content.Flavor, id content.FlavorId) int {
	for i, f := range flavors {
		if f.Id == id {
			return i
		}
	}
	return 0
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
