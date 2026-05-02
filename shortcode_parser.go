package demoji

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// shortcodeInlineParser converts :shortcode: patterns to unicode emoji.
// It must be an inline parser (not an AST transformer) because shortcode
// names can contain underscores, and goldmark's emphasis tokenizer fragments
// text at '_' boundaries before any AST transformer runs — a transformer
// would never see ":slightly_smiling_face:" as a single string.
//
// By triggering on ':', we consume the complete :name: token before the
// emphasis parser can look inside it.
type shortcodeInlineParser struct {
	table map[string]string // :shortcode: → emoji
}

func (p *shortcodeInlineParser) Trigger() []byte {
	return []byte{':'}
}

func (p *shortcodeInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) < 3 || line[0] != ':' {
		return nil
	}
	i := 1
	for i < len(line) && isShortcodeChar(line[i]) {
		i++
	}
	if i < 2 || i >= len(line) || line[i] != ':' {
		return nil
	}
	emoji, ok := p.table[string(line[:i+1])]
	if !ok {
		return nil
	}
	block.Advance(i + 1)
	return ast.NewString([]byte(emoji))
}

func isShortcodeChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= '0' && b <= '9') ||
		b == '_' || b == '-' || b == '+'
}
