package edit

import (
	"io"

	"github.com/polyglottis/frontend_server/extract"
	"github.com/polyglottis/frontend_server/server"
	"github.com/polyglottis/frontend_server/templates"
	"github.com/polyglottis/platform/content"
	"github.com/polyglottis/platform/frontend"
	"github.com/polyglottis/platform/i18n"
)

type EditServer struct{}

var editTextTmpl = templates.Parse("extract/templates/actions.html", "extract/edit/templates/text.html")

func (s *EditServer) EditText(context *frontend.Context, e *content.Extract, a, b *content.Flavor) ([]byte, error) {
	return server.Call(context, func(w io.Writer, serverArgs *server.TmplArgs) error {
		args := &TmplArgs{
			TmplArgs:     extract.NewTmplArgs(serverArgs, e, a, b),
			ExtractShape: e.Shape(),
		}
		if context.IsFocusOnA() {
			args.Focus, args.NoFocus = a, b
		} else {
			args.Focus, args.NoFocus = b, a
		}
		defaults := context.Defaults
		if sum := defaults.Get("Summary"); sum == "" {
			defaults.Set("Summary", args.Focus.Summary)
		}
		if title := defaults.Get("Title"); title == "" {
			defaults.Set("Title", args.TitleUnderFocus())
		}
		args.Data = map[string]interface{}{
			"title":    getTitle(args.Focus),
			"HasA":     a != nil,
			"HasB":     b != nil,
			"errors":   context.Errors,
			"defaults": defaults,
		}
		args.Css = "edit"
		return editTextTmpl.Execute(w, args)
	})
}

type TmplArgs struct {
	*extract.TmplArgs
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
	if a.Context.IsFocusOnA() {
		return a.LanguageA()
	} else {
		return a.LanguageB()
	}
}

func (a *TmplArgs) LanguageOther() string {
	if a.Context.IsFocusOnA() {
		return a.LanguageB()
	} else {
		return a.LanguageA()
	}
}

func (a *TmplArgs) OtherLanguage() bool {
	return a.NoFocus != nil
}

func (a *TmplArgs) TitleUnderFocus() string {
	title := a.Focus.GetTitle()
	if title == nil {
		return ""
	}
	return title.Content
}

func (a *TmplArgs) TitleOther() string {
	title := a.NoFocus.GetTitle()
	if title == nil {
		return ""
	}
	return title.Content
}

type editBlock struct {
	BlockId content.BlockId
	Lines   []*editLine
}

type editLine struct {
	UnitId            content.UnitId
	ContentUnderFocus string
	ContentOther      string
}

func (a *TmplArgs) EditBlocks() []*editBlock {
	blocks := make([]*editBlock, 0, len(a.ExtractShape))
	a.ExtractShape.IterateFlavorBody(a.Focus, func(blockId content.BlockId) {
		blocks = append(blocks, &editBlock{
			BlockId: blockId,
		})
	}, func(blockId content.BlockId, unitId content.UnitId, u *content.Unit) {
		line := &editLine{UnitId: unitId}
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
