package minify

import (
	"io"
	"regexp"

	minifypkg "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-scriptlet/render"
)

var minify *minifypkg.M

func init() {
	minify = minifypkg.New()

	minify.AddFunc("text/css", css.Minify)

	minify.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
	})

	minify.AddFunc("image/svg+xml", svg.Minify)

	minify.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	minify.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

	minify.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

// ([scriptlet.RenderFunc] signature)
func RenderMinifyCSS(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMinify("text/css", writer, content, js)
}

// ([scriptlet.RenderFunc] signature)
func RenderMinifyHTML(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMinify("text/html", writer, content, js)
}

// ([scriptlet.RenderFunc] signature)
func RenderMinifySVG(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMinify("image/svg+xml", writer, content, js)
}

// ([scriptlet.RenderFunc] signature)
func RenderMinifyJavaScript(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMinify("text/javascript", writer, content, js)
}

// ([scriptlet.RenderFunc] signature)
func RenderMinifyJSON(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMinify("application/json", writer, content, js)
}

// ([scriptlet.RenderFunc] signature)
func RenderMinifyXML(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	return renderMinify("application/xml", writer, content, js)
}

// Utils

func renderMinify(mediaType string, writer io.Writer, content any, js bool) error {
	reader := render.ToReader(content)
	if js {
		return render.AsPresenter(writer, minify.Reader(mediaType, reader))
	} else {
		return minify.Minify(mediaType, writer, reader)
	}
}
