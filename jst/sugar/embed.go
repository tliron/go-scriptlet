package sugar

import (
	"strings"

	"github.com/tliron/go-scriptlet/jst"
)

// Example:
//
//	this.clone().embed(bind(CODE, 'present'), env.context);
//
// ([jst.HandleSugarFunc] signature)
func HandleEmbed(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):]

	if err := scriptletContext.WriteString("this.clone().embed(bind("); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
		return false, err
	}
	return false, scriptletContext.WriteString(", 'present'), env.context);\n")
}
