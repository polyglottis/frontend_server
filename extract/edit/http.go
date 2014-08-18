package edit

import (
	"io"

	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type EditServer struct{}

var editTextTmpl = templates.Parse("extract/templates/frame", "extract/edit/templates/text")

func (s *EditServer) EditText(context *frontend.Context, extract *content.Extract, a, b *content.Flavor) ([]byte, error) {
	return server.Call(context, func(w io.Writer, serverArgs *server.TmplArgs) error {
		args := &TmplArgs{
			TmplArgs:     serverArgs,
			ExtractShape: extract.Shape(),
		}
		if context.Query.Get("focus") != "b" {
			context.Query.Set("focus", "a")
			args.Focus, args.NoFocus = a, b
		} else {
			args.Focus, args.NoFocus = b, a
		}
		args.Data = map[string]interface{}{
			"title": getTitle(args.Focus),
		}
		return editTextTmpl.Execute(w, args)
	})
}

type TmplArgs struct {
	*server.TmplArgs
	ExtractShape content.ExtractShape
	Focus        *content.Flavor
	NoFocus      *content.Flavor
}

func getTitle(f *content.Flavor) interface{} {
	if len(f.Blocks) == 0 || f.Blocks[0][0].BlockId != 1 {
		return i18n.Key("No title")
	} else {
		return f.Blocks[0][0].Content
	}
}

func (a *TmplArgs) LanguageUnderFocus() string {
	if a.Context.Query.Get("focus") == "a" {
		return a.LanguageA()
	} else {
		return a.LanguageB()
	}
}

func (a *TmplArgs) LanguageOther() string {
	if a.Context.Query.Get("focus") == "a" {
		return a.LanguageB()
	} else {
		return a.LanguageA()
	}
}

func (a *TmplArgs) OtherLanguage() bool {
	return a.NoFocus != nil
}

func (a *TmplArgs) TitleUnderFocus() string {
	f := a.Focus
	if len(f.Blocks) == 0 || f.Blocks[0][0].BlockId != 1 {
		return a.GetText(i18n.Key("Type the title here"))
	}
	return f.Blocks[0][0].Content
}

func (a *TmplArgs) TitleOther() string {
	f := a.NoFocus
	if f == nil || len(f.Blocks) == 0 || f.Blocks[0][0].BlockId != 1 {
		return ""
	}
	return f.Blocks[0][0].Content
}

type editBlock struct {
	Id    content.BlockId
	Lines []*editLine
}

type editLine struct {
	ContentUnderFocus string
	ContentOther      string
}

func (a *TmplArgs) EditBlocks() []*editBlock {
	blocks := make([]*editBlock, 0, len(a.ExtractShape))
	a.ExtractShape.IterateFlavorBody(a.Focus, func(blockId content.BlockId) {
		blocks = append(blocks, &editBlock{
			Id: blockId,
		})
	}, func(blockId content.BlockId, unitId content.UnitId, u *content.Unit) {
		line := &editLine{}
		if u != nil {
			line.ContentUnderFocus = u.Content
		}
		blocks[len(blocks)-1].Lines = append(blocks[len(blocks)-1].Lines, line)
	}, nil)
	if a.NoFocus != nil {
		a.ExtractShape.IterateFlavorBody(a.NoFocus, nil, func(blockId content.BlockId, unitId content.UnitId, u *content.Unit) {
			if u != nil {
				blocks[int(blockId)-2].Lines[int(unitId)-1].ContentOther = u.Content
			}
		}, nil)
	}
	return blocks
}
