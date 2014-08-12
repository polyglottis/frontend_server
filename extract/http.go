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

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

var flavorTmpl = templates.Parse("extract/templates/flavor")

func (s *Server) Extract(context *frontend.Context, extract *content.Extract) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		f := extractFlavor(extract.Flavors[0])
		args.Data = map[string]interface{}{
			"title":   getTitle(f),
			"flavorA": f,
		}
		return flavorTmpl.Execute(w, args)
	})
}

func (s *Server) Flavor(context *frontend.Context, extract *content.Extract, flavor *content.Flavor) ([]byte, error) {
	return server.Call(context, func(w io.Writer, args *server.TmplArgs) error {
		f := extractFlavor(flavor)
		args.Data = map[string]interface{}{
			"title":   getTitle(f),
			"flavorA": f,
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

func extractFlavor(flavor *content.Flavor) *Flavor {
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
	f.Blocks = extractBlocks(flavor.Blocks)
	return f
}

func extractBlocks(blocks content.BlockSlice) []*Block {
	if len(blocks) == 0 {
		return []*Block{{Missing: true}}
	}
	b := make([]*Block, 0, len(blocks))
	lastBlockId := content.BlockId(1)
	for _, block := range blocks {
		curBlockId := block[0].BlockId
		if curBlockId == 1 {
			continue
		}
		for ; lastBlockId < curBlockId-1; lastBlockId++ {
			b = append(b, &Block{Id: lastBlockId, Missing: true})
		}
		lastBlockId = curBlockId
		b = append(b, &Block{Id: curBlockId, Units: extractUnits(block)})
	}
	if len(b) == 0 {
		return []*Block{{Missing: true}}
	}
	return b
}

func extractUnits(units content.UnitSlice) []*Unit {
	u := make([]*Unit, 0, len(units))
	lastUnitId := content.UnitId(0)
	for _, unit := range units {
		for ; lastUnitId < unit.Id-1; lastUnitId++ {
			u = append(u, &Unit{Id: lastUnitId, String: String{Missing: true}})
		}
		lastUnitId = unit.Id
		u = append(u, &Unit{Id: unit.Id, String: String{Content: unit.Content}})
	}
	return u
}
