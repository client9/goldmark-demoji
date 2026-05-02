# goldmark-demoji

A [goldmark](https://github.com/yuin/goldmark) extension that converts between
Unicode emoji, ASCII emoticons, and GitHub-style shortcodes in Markdown documents.

```
😀  ↔  :-D                       (unicode ↔ emoticon)
🙂  ↔  :slightly_smiling_face:   (unicode ↔ shortcode)
```

Emoji inside code spans and fenced code blocks are left unchanged.

For plain-string conversion without goldmark, use
[github.com/client9/demoji](https://github.com/client9/demoji) directly.

## Install

```bash
go get github.com/client9/goldmark-demoji
```

## Quick start

```go
import (
    demoji "github.com/client9/goldmark-demoji"
    "github.com/yuin/goldmark"
)

// Default: unicode emoji → ASCII emoticons.
md := goldmark.New(goldmark.WithExtensions(demoji.New()))
```

## Formats

| Constant | Example | Description |
|----------|---------|-------------|
| `FormatUnicode` | `😀` | Unicode emoji codepoints |
| `FormatEmoticon` | `:-D` | Classic ASCII emoticons |
| `FormatShortcode` | `:grinning:` | GitHub / Gemoji shortcodes |

## Conversion directions

```go
// unicode → emoticon (default)
demoji.New()

// unicode → shortcode
demoji.New(
    demoji.WithFrom(demoji.FormatUnicode),
    demoji.WithTo(demoji.FormatShortcode),
)

// emoticon → unicode
demoji.New(demoji.WithFrom(demoji.FormatEmoticon))

// shortcode → unicode
demoji.New(demoji.WithFrom(demoji.FormatShortcode))

// emoticon + shortcode → unicode (single pass)
demoji.New(demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode))
```

## Configuration

### WithAdditional

Add or override entries. Keys are source patterns; values are Unicode emoji strings.

```go
// Add a custom emoticon.
demoji.New(
    demoji.WithFrom(demoji.FormatEmoticon),
    demoji.WithAdditional(demoji.FormatEmoticon, map[string]string{"^^": "😊"}),
)

// Add a custom shortcode.
demoji.New(
    demoji.WithFrom(demoji.FormatShortcode),
    demoji.WithAdditional(demoji.FormatShortcode, map[string]string{":unicorn:": "🦄"}),
)
```

### WithExclude

Remove entries from the built-in mapping.

```go
// Prevent 🙂 from being converted at all (both canonical forms excluded).
demoji.New(
    demoji.WithExclude(demoji.FormatEmoticon, ":-)", ":)"),
)
```

### Canonical text selection

In unicode→text mode, each emoji maps to one canonical text form — the first
active entry in the built-in table. Exclude the current leader to shift the
canonical to the next entry:

```go
// default: 🙂 → :-)
// after removing ":-)", next entry ":)" becomes canonical:
demoji.New(
    demoji.WithExclude(demoji.FormatEmoticon, ":-)"),
)
// → 🙂 renders as :)
```

## Code preservation

Emoji and emoticons inside code spans and fenced code blocks are never modified:

```markdown
outside 🙂 but `inside 🙂 code span`
```

renders as:

```html
<p>outside :-) but <code>inside 🙂 code span</code></p>
```
