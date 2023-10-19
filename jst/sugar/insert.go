package sugar

import (
	"strings"

	"github.com/tliron/go-scriptlet/jst"
)

// Example:
//
//	const __args0 = [CODE];
//	if (__args0.length>1) scriptlet.renderFrom(this.writer, __args0[0], __args0[1]);
//	else env.writeFrom(this.writer, __args0[0]);
//
// ([jst.HandleSugarFunc] signature)
func HandleInsert(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):]
	suffix := scriptletContext.NextSuffix()

	if err := scriptletContext.WriteString("const __args"); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(suffix); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(" = ["); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(strings.TrimSpace(code)); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString("];\n"); err != nil {
		return false, err
	}

	if err := scriptletContext.WriteString("if (__args"); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(suffix); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(".length>1) scriptlet.renderFrom(this.writer, __args"); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(suffix); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString("[0], __args"); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(suffix); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString("[1]);\n"); err != nil {
		return false, err
	}

	if err := scriptletContext.WriteString("else env.writeFrom(this.writer, __args"); err != nil {
		return false, err
	}
	if err := scriptletContext.WriteString(suffix); err != nil {
		return false, err
	}
	return false, scriptletContext.WriteString("[0]);\n")
}
