package minify

import (
	"github.com/tliron/go-scriptlet/render"
)

func RegisterDefaultRenderers() {
	render.RegisterRenderer("mincss", RenderMinifyCSS)
	render.RegisterRenderer("minhtml", RenderMinifyHTML)
	render.RegisterRenderer("minsvg", RenderMinifySVG)
	render.RegisterRenderer("minjs", RenderMinifyJavaScript)
	render.RegisterRenderer("minjson", RenderMinifyJSON)
	render.RegisterRenderer("minxml", RenderMinifyXML)
}

func RegisterAutomaticRenderers() {
	render.RegisterRenderer("css", RenderMinifyCSS)
	render.RegisterRenderer("html", RenderMinifyHTML)
	render.RegisterRenderer("svg", RenderMinifySVG)
	render.RegisterRenderer("js", RenderMinifyJavaScript)
	render.RegisterRenderer("json", RenderMinifyJSON)
	render.RegisterRenderer("xml", RenderMinifyXML)
}
