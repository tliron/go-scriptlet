package sugar

import (
	"strings"

	"github.com/tliron/go-scriptlet/jst"
)

// Examples:
//
//	this.startRender(CODE, env.context);
//
//	this.endRender();
//
// ([jst.HandleSugarFunc] signature)
func HandleRender(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):]

	if code == prefix {
		// End render
		return false, scriptletContext.WriteString("this.endRender();\n")
	} else {
		// Start render
		if err := scriptletContext.WriteString("this.startRender("); err != nil {
			return false, err
		}
		if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
			return false, err
		}
		return false, scriptletContext.WriteString(", env.context);\n")
	}
}
