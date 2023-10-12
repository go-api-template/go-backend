package mailer

import (
	"github.com/gomarkdown/markdown"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

// DotEscape escapes dots in a string
// It wraps a dots in names with ZERO WIDTH JOINER [U+200D] in order to prevent
// autolinkers from detecting these as urls
func DotEscape(raw string) string {
	return strings.ReplaceAll(raw, ".", "\u200d.\u200d")
}

// Str2html render Markdown text to HTML
func Str2html(raw string) string {
	unsafeHtml := markdown.ToHTML([]byte(raw), nil, nil)
	return string(bluemonday.UGCPolicy().SanitizeBytes(unsafeHtml))
}
