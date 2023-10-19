package sugar

import (
	"strings"

	"github.com/tliron/go-scriptlet/jst"
)

const htmlEscapePrefix = "/"

// Examples:
//
//	this.write(CODE);
//
//	this.write(util.escapeHtml(CODE));
//
//	this.write(this.getVariable(CODE));
//
//	this.write(util.escapeHtml(this.getVariable(CODE)));
//
// ([jst.HandleSugarFunc] signature)
func HandleExpression(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	prefixLength := len(prefix)
	code = code[prefixLength:]

	if strings.HasPrefix(code, prefix) {
		code = code[prefixLength:]

		if strings.HasPrefix(code, htmlEscapePrefix) {
			code = code[len(htmlEscapePrefix):]

			if err := scriptletContext.WriteString("this.write(util.escapeHtml(this.getVariable("); err != nil {
				return true, err
			}
			if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
				return true, err
			}
			return true, scriptletContext.WriteString(")));\n")
		} else {
			if err := scriptletContext.WriteString("this.write(this.getVariable("); err != nil {
				return true, err
			}
			if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
				return true, err
			}
			return true, scriptletContext.WriteString("));\n")
		}
	} else {
		if strings.HasPrefix(code, htmlEscapePrefix) {
			code = code[len(htmlEscapePrefix):]

			if err := scriptletContext.WriteString("this.write(util.escapeHtml("); err != nil {
				return true, err
			}
			if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
				return true, err
			}
			return true, scriptletContext.WriteString("));\n")
		} else {
			if err := scriptletContext.WriteString("this.write("); err != nil {
				return true, err
			}
			if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
				return true, err
			}
			return true, scriptletContext.WriteString(");\n")
		}
	}
}
