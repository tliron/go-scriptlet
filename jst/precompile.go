package jst

import (
	"io"
	"strings"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/exturl"
	"github.com/tliron/go-scriptlet/render"
)

var DebugWriter io.Writer

// ([commonjs.PrecompileFunc] signature)
func Precompile(url exturl.URL, script string, jsContext *commonjs.Context) (string, error) {
	format := url.Format()
	for _, renderer := range render.GetRenderers() {
		if format == renderer {
			content, err := render.RenderToString(script, format, true, jsContext)

			if (DebugWriter != nil) && (err == nil) {
				io.WriteString(DebugWriter, strings.TrimRight(content, "\n"))
			}

			return content, err
		}
	}

	return script, nil
}
