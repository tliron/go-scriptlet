package sanitize

import (
	"github.com/tliron/go-scriptlet/render"
)

func RegisterDefaultRenderers() {
	render.RegisterRenderer("sanitizehtml", RenderSanitizeHTML)
}
