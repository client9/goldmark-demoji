// Package demoji provides a goldmark extension that converts between Unicode
// emoji, ASCII emoticons, and GitHub-style shortcodes.
//
// All Format constants and Option constructors are defined here so callers
// need only one import. For plain-string conversion without goldmark, use
// github.com/client9/demoji directly.
package demoji

import (
	"maps"

	dj "github.com/client9/demoji"
	gm "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

// Format is an alias for the core Format type.
type Format = dj.Format

// Format constants re-exported by value.
const (
	FormatUnicode   = dj.FormatUnicode
	FormatEmoticon  = dj.FormatEmoticon
	FormatShortcode = dj.FormatShortcode
)

// Option configures the extension.
type Option func(*dj.Config)

// WithFrom sets the source format(s). OR multiple values to convert several
// formats in one AST pass:
//
//	demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode)
func WithFrom(f Format) Option {
	return func(cfg *dj.Config) { cfg.From = f }
}

// WithTo sets the target format. Only used when From includes FormatUnicode.
func WithTo(f Format) Option {
	return func(cfg *dj.Config) { cfg.To = f }
}

// WithAdditional adds or overrides entries in the mapping for format.
// Keys are source patterns; values are unicode emoji strings.
func WithAdditional(format Format, m map[string]string) Option {
	return func(cfg *dj.Config) {
		switch format {
		case dj.FormatEmoticon:
			if cfg.Emoticons == nil {
				cfg.Emoticons = maps.Clone(dj.DefaultEmoticons())
			}
			maps.Copy(cfg.Emoticons, m)
		case dj.FormatShortcode:
			if cfg.Shortcodes == nil {
				cfg.Shortcodes = maps.Clone(dj.DefaultShortcodes())
			}
			maps.Copy(cfg.Shortcodes, m)
		}
	}
}

// WithExclude removes the named patterns from the mapping for format.
func WithExclude(format Format, keys ...string) Option {
	return func(cfg *dj.Config) {
		switch format {
		case dj.FormatEmoticon:
			if cfg.Emoticons == nil {
				cfg.Emoticons = maps.Clone(dj.DefaultEmoticons())
			}
			for _, k := range keys {
				delete(cfg.Emoticons, k)
			}
		case dj.FormatShortcode:
			if cfg.Shortcodes == nil {
				cfg.Shortcodes = maps.Clone(dj.DefaultShortcodes())
			}
			for _, k := range keys {
				delete(cfg.Shortcodes, k)
			}
		}
	}
}

// Extension is the goldmark extension. Create with New.
type Extension struct {
	r          *dj.Replacer
	shortcodes map[string]string
	from       dj.Format
}

// New creates the extension. With no options the default behavior is
// unicode emoji → ASCII emoticons.
func New(opts ...Option) *Extension {
	cfg := dj.DefaultConfig()
	for _, o := range opts {
		o(&cfg)
	}

	// Resolve the shortcodes map now so the inline parser has a concrete reference.
	shortcodes := cfg.Shortcodes
	if shortcodes == nil {
		shortcodes = dj.DefaultShortcodes()
	}

	// The AST transformer must not process shortcodes → unicode: goldmark's
	// emphasis tokenizer fragments text at '_' boundaries before any transformer
	// runs, so ":slightly_smiling_face:" would never appear as a single string.
	// The inline parser handles that direction; strip FormatShortcode from the
	// transformer's config so its strings.Replacer omits those pairs.
	transformerCfg := cfg
	transformerCfg.From &^= dj.FormatShortcode

	return &Extension{
		r:          dj.New(transformerCfg),
		shortcodes: shortcodes,
		from:       cfg.From,
	}
}

// Extend implements goldmark.Extender.
func (e *Extension) Extend(m gm.Markdown) {
	// Register the inline parser only when converting shortcodes → unicode.
	if e.from&dj.FormatShortcode != 0 && e.from&dj.FormatUnicode == 0 {
		m.Parser().AddOptions(
			parser.WithInlineParsers(
				util.Prioritized(&shortcodeInlineParser{table: e.shortcodes}, 999),
			),
		)
	}
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(&emojiTransformer{r: e.r}, 100),
		),
	)
}
