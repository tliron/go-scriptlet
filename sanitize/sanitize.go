package sanitize

import (
	"io"

	"github.com/microcosm-cc/bluemonday"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-scriptlet/render"
)

var policy = bluemonday.UGCPolicy()

// ([scriptlet.RenderFunc] signature)
func RenderSanitizeHTML(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	if content_, err := render.ToBytes(content); err == nil {
		content_ = policy.SanitizeBytes(content_)
		if js {
			return render.AsPresenter(writer, content_)
		} else {
			_, err := render.Write(writer, content_)
			return err
		}
	} else {
		return err
	}
}
