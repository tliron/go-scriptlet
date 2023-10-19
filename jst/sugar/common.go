package sugar

import (
	"github.com/tliron/go-scriptlet/jst"
)

func RegisterDefaultSugar() {
	jst.RegisterSugar("#", HandleComment)
	jst.RegisterSugar("&", HandleEmbed)
	jst.RegisterSugar("=", HandleExpression)
	jst.RegisterSugar("+", HandleInsert)
	jst.RegisterSugar("!", HandleCapture)
	jst.RegisterSugar("^", HandleRender)
}
