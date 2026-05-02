package demoji_test

import (
	"bytes"
	"testing"

	demoji "github.com/client9/goldmark-demoji"
	gm "github.com/yuin/goldmark"
)

func render(t *testing.T, ext *demoji.Extension, src string) string {
	t.Helper()
	md := gm.New(gm.WithExtensions(ext))
	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		t.Fatalf("Convert: %v", err)
	}
	return buf.String()
}

// --- unicode → emoticon (default) ---

func TestDemojiToEmoticon(t *testing.T) {
	ext := demoji.New()
	tests := []struct{ name, in, want string }{
		{"smile", "Hello 🙂 world", "<p>Hello :-) world</p>\n"},
		{"grin", "Woohoo 😀", "<p>Woohoo :-D</p>\n"},
		{"multiple", "😀 and 😉", "<p>:-D and ;-)</p>\n"},
		{"love", "I ❤️ Go", "<p>I &lt;3 Go</p>\n"},
		{"broken heart", "💔", "<p>&lt;/3</p>\n"},
		{"no emoji", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- unicode → shortcode ---

func TestDemojiToShortcode(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatUnicode),
		demoji.WithTo(demoji.FormatShortcode),
	)
	tests := []struct{ name, in, want string }{
		{"smile", "Hello 🙂 world", "<p>Hello :slightly_smiling_face: world</p>\n"},
		{"grin", "Woohoo 😀", "<p>Woohoo :grinning:</p>\n"},
		{"heart", "I ❤️ Go", "<p>I :heart: Go</p>\n"},
		{"no emoji", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- emoticon → unicode ---

func TestRemojiFromEmoticon(t *testing.T) {
	ext := demoji.New(demoji.WithFrom(demoji.FormatEmoticon))
	tests := []struct{ name, in, want string }{
		{"classic smile", "Hello :-) world", "<p>Hello 🙂 world</p>\n"},
		{"short smile", "Hi :)", "<p>Hi 🙂</p>\n"},
		{"wink", ";-) wink", "<p>😉 wink</p>\n"},
		{"angel longer first", "O:-) angel", "<p>😇 angel</p>\n"},
		{"love", "I <3 Go", "<p>I ❤️ Go</p>\n"},
		{"multiple", ":-) and :-(", "<p>🙂 and 🙁</p>\n"},
		{"no emoticons", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- shortcode → unicode ---

func TestRemojiFromShortcode(t *testing.T) {
	ext := demoji.New(demoji.WithFrom(demoji.FormatShortcode))
	tests := []struct{ name, in, want string }{
		{"smile", "Hello :slightly_smiling_face: world", "<p>Hello 🙂 world</p>\n"},
		{"grin", ":grinning:", "<p>😀</p>\n"},
		{"long shortcode first", ":stuck_out_tongue_winking_eye:", "<p>😜</p>\n"},
		{"heart", "I :heart: Go", "<p>I ❤ Go</p>\n"},
		{"no shortcodes", "Plain text", "<p>Plain text</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- emoticon | shortcode → unicode (single pass, bitfield) ---

func TestRemojiCombined(t *testing.T) {
	ext := demoji.New(demoji.WithFrom(demoji.FormatEmoticon | demoji.FormatShortcode))
	tests := []struct{ name, in, want string }{
		{
			"emoticon and shortcode",
			":-) and :grinning: are both happy",
			"<p>🙂 and 😀 are both happy</p>\n",
		},
		{"shortcode only", ":wink:", "<p>😉</p>\n"},
		{"emoticon only", ";-)", "<p>😉</p>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := render(t, ext, tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

// --- code blocks and spans are skipped ---

func TestSkipsCodeSpan(t *testing.T) {
	tests := []struct {
		mode demoji.Format
		in   string
		want string
	}{
		{
			demoji.FormatUnicode,
			"outside 🙂 but `inside 🙂 code`",
			"<p>outside :-) but <code>inside 🙂 code</code></p>\n",
		},
		{
			demoji.FormatEmoticon,
			"outside :-) but `inside :-)  code`",
			"<p>outside 🙂 but <code>inside :-)  code</code></p>\n",
		},
		{
			demoji.FormatShortcode,
			"outside :grinning: but `inside :grinning: code`",
			"<p>outside 😀 but <code>inside :grinning: code</code></p>\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.in[:15], func(t *testing.T) {
			got := render(t, demoji.New(demoji.WithFrom(tt.mode)), tt.in)
			if got != tt.want {
				t.Errorf("got  %q\nwant %q", got, tt.want)
			}
		})
	}
}

func TestSkipsFencedCodeBlock(t *testing.T) {
	got := render(t, demoji.New(), "```\n🙂 :-) :grinning:\n```")
	want := "<pre><code>🙂 :-) :grinning:\n</code></pre>\n"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

// --- configuration: WithAdditional and WithExclude ---

func TestWithAdditional(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatEmoticon),
		demoji.WithAdditional(demoji.FormatEmoticon, map[string]string{"^^": "😊"}),
	)
	got := render(t, ext, "Hello ^^ world")
	if got != "<p>Hello 😊 world</p>\n" {
		t.Errorf("got %q", got)
	}
}

func TestWithAdditionalShortcode(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatShortcode),
		demoji.WithAdditional(demoji.FormatShortcode, map[string]string{":unicorn:": "🦄"}),
	)
	got := render(t, ext, "Hello :unicorn: world")
	if got != "<p>Hello 🦄 world</p>\n" {
		t.Errorf("got %q", got)
	}
}

func TestWithExclude(t *testing.T) {
	ext := demoji.New(
		demoji.WithExclude(demoji.FormatEmoticon, ":-)", ":)"),
	)
	got := render(t, ext, "Hello 🙂 world")
	if got != "<p>Hello 🙂 world</p>\n" {
		t.Errorf("got %q", got)
	}
}

func TestWithExcludeCanonicalFallback(t *testing.T) {
	ext := demoji.New(
		demoji.WithExclude(demoji.FormatEmoticon, ":-)"),
	)
	got := render(t, ext, "Hello 🙂 world")
	if got != "<p>Hello :) world</p>\n" {
		t.Errorf("got %q", got)
	}
}

func TestWithExcludeShortcode(t *testing.T) {
	ext := demoji.New(
		demoji.WithFrom(demoji.FormatUnicode),
		demoji.WithTo(demoji.FormatShortcode),
		demoji.WithExclude(demoji.FormatShortcode, ":grinning:"),
	)
	got := render(t, ext, "Woohoo 😀")
	if got != "<p>Woohoo 😀</p>\n" {
		t.Errorf("got %q", got)
	}
}

func BenchmarkNew(b *testing.B) {
	b.Run("no_options", func(b *testing.B) {
		for range b.N {
			_ = demoji.New()
		}
	})
	b.Run("with_additional", func(b *testing.B) {
		extra := map[string]string{"^^": "😊"}
		for range b.N {
			_ = demoji.New(demoji.WithAdditional(demoji.FormatEmoticon, extra))
		}
	})
}
