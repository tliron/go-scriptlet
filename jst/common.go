package jst

import (
	"github.com/tliron/go-scriptlet/render"
)

func RegisterDefaultRenderers() {
	render.RegisterRenderer("jst", RenderJST)
}
