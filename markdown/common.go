package markdown

import (
	"github.com/tliron/go-scriptlet/render"
)

func RegisterDefaultRenderers() {
	render.RegisterRenderer("markdown", RenderMarkdown)
	render.RegisterRenderer("md", RenderMarkdown)
	render.RegisterRenderer("extendedmarkdown", RenderExtendedMarkdown)
	render.RegisterRenderer("extendedmd", RenderExtendedMarkdown)
}
