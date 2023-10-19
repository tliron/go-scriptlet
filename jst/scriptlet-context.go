package jst

import (
	"io"
	"strconv"

	"github.com/tliron/go-scriptlet/render"
)

//
// ScriptletContext
//

type ScriptletContext struct {
	Writer io.Writer

	embedIndex int64
}

func NewScriptletContext(writer io.Writer) *ScriptletContext {
	return &ScriptletContext{Writer: writer}
}

func (self *ScriptletContext) NextSuffix() string {
	suffix := strconv.FormatInt(self.embedIndex, 10)
	self.embedIndex++
	return suffix
}

func (self *ScriptletContext) WriteString(content string) error {
	if content == "" {
		return nil
	}

	_, err := render.WriteString(self.Writer, content)
	return err
}

func (self *ScriptletContext) AsContextWrite(literal string) error {
	if literal == "" {
		return nil
	}

	return render.AsContextWrite(self.Writer, literal)
}
