package markdown

import (
	"bytes"
	"io"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-scriptlet/render"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var commonMarkdown = goldmark.New(
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

var extendedMarkdown = goldmark.New(
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
		parser.WithAttribute(),
	),
	goldmark.WithExtensions(
		extension.GFM,            // https://github.github.com/gfm/
		extension.DefinitionList, // https://michelf.ca/projects/php-markdown/extra/#def-list
		extension.Footnote,       // https://michelf.ca/projects/php-markdown/extra/#footnotes
		extension.Typographer,    // https://daringfireball.net/projects/smartypants/
		extension.CJK,
	),
)

// ([scriptlet.RenderFunc] signature)
func RenderMarkdown(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMarkdown(writer, content, js, commonMarkdown)
}

// ([scriptlet.RenderFunc] signature)
func RenderExtendedMarkdown(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMarkdown(writer, content, js, extendedMarkdown)
}

func renderMarkdown(writer io.Writer, content any, js bool, markdown goldmark.Markdown) error {
	if content_, err := render.ToBytes(content); err == nil {
		if js {
			var buffer bytes.Buffer
			if err := markdown.Convert(content_, &buffer); err == nil {
				return render.AsPresenter(writer, buffer.Bytes())
			} else {
				return err
			}
		} else {
			return markdown.Convert(content_, writer)
		}
	} else {
		return err
	}
}
