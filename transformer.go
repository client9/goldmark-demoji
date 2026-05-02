package demoji

import (
	dj "github.com/client9/demoji"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type emojiTransformer struct {
	r *dj.Replacer
}

func (t *emojiTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	source := reader.Source()
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n.Kind() {
		case ast.KindCodeBlock, ast.KindFencedCodeBlock, ast.KindCodeSpan,
			ast.KindHTMLBlock, ast.KindRawHTML:
			return ast.WalkSkipChildren, nil
		case ast.KindText:
			replaceTextNode(n.(*ast.Text), source, t.r)
		}
		return ast.WalkContinue, nil
	})
}

func replaceTextNode(node *ast.Text, source []byte, r *dj.Replacer) {
	src := string(node.Segment.Value(source))
	result := r.Replace(src)
	if result == src {
		return
	}
	node.Parent().ReplaceChild(node.Parent(), node, ast.NewString([]byte(result)))
}
