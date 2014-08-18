package extract

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
	"github.com/polyglottis/platform/language"
)

type ExtractServer struct{}

var flavorTmpl = templates.Parse("extract/templates/frame", "extract/templates/flavor")

func (s *ExtractServer) Extract(context *frontend.Context, extract *content.Extract) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		args.Data = map[string]interface{}{
			"title": "home_page",
		}
		return flavorTmpl.Execute(w, args)
	})
}

func (s *ExtractServer) Flavor(context *frontend.Context, extract *content.Extract, a, b *frontend.FlavorTriple) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		f := extractFlavor(extract.Shape(), a.Text)
		args.Data = map[string]interface{}{
			"title":     getTitle(f),
			"flavorA":   f,
			"languageA": context.LanguageA,
		}
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

func extractFlavor(shape content.ExtractShape, flavor *content.Flavor) *Flavor {
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
