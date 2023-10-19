package sugar

import (
	"strings"

	"github.com/tliron/go-scriptlet/jst"
)

// Examples:
//
//	this.startCapture(CODE);
//
//	this.endCapture();
//
// ([jst.HandleSugarFunc] signature)
func HandleCapture(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):]

	if code == prefix {
		// End render
		return false, scriptletContext.WriteString("this.endCapture();\n")
	} else {
		// Start render
		if err := scriptletContext.WriteString("this.startCapture("); err != nil {
			return false, err
		}
		if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
			return false, err
		}
		return false, scriptletContext.WriteString(");\n")
	}
}
