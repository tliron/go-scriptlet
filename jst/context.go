package jst

import (
	"errors"
	"io"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-scriptlet/render"
)

//
// Context
//

type Context struct {
	Variables ard.StringMap
	Writer    io.Writer
}

// Note: the expectation is that the writer is not asynchronous.
// See [WriteString].
func NewContext(writer io.Writer, variables ard.StringMap) *Context {
	if variables == nil {
		variables = make(ard.StringMap)
	}
	return &Context{
		Variables: variables,
		Writer:    writer,
	}
}

func (self *Context) Clone() *Context {
	return &Context{
		Variables: ard.Copy(self.Variables).(ard.StringMap),
		Writer:    self.Writer,
	}
}

func (self *Context) GetVariable(keys ...any) any {
	if value := ard.With(self.Variables).ConvertSimilar().Get(keys...).Value; value != nil {
		return value
	} else {
		return goja.Undefined()
	}
}

func (self *Context) Write(content any) error {
	_, err := render.Write(self.Writer, content)
	return err
}

func (self *Context) StartCapture(name string) {
	self.Writer = NewCaptureWriter(self.Writer, name, func(name string, value string) {
		self.Variables[name] = value
	})
}

func (self *Context) EndCapture() error {
	if captureWriter, ok := self.Writer.(*CaptureWriter); ok {
		err := captureWriter.Close()
		self.Writer = captureWriter.originalWriter
		return err
	} else {
		return errors.New("did not call startCapture()")
	}
}

func (self *Context) StartRender(renderer string, jsContext *commonjs.Context) error {
	if renderWriter, err := NewRenderWriter(self.Writer, jsContext, renderer); err == nil {
		self.Writer = renderWriter
		return nil
	} else {
		return err
	}
}

func (self *Context) EndRender() error {
	if renderWriter, ok := self.Writer.(*RenderWriter); ok {
		err := renderWriter.Close()
		self.Writer = renderWriter.originalWriter
		return err
	} else {
		return errors.New("did not call startRender()")
	}
}

func (self *Context) Embed(present any, jsContext *commonjs.Context) error {
	var err error
	if present, jsContext, err = commonjs.Unbind(present, jsContext); err != nil {
		return err
	}

	_, err = jsContext.Environment.Call(present, self)
	return err
}

func (self *Context) Flush() error {
	var err error
	if self.Writer, err = UnwrapWriters(self.Writer); err == nil {
		if closer, ok := self.Writer.(io.Closer); ok {
			return closer.Close()
		} else {
			return nil
		}
	} else {
		return err
	}
}
