package jst

import (
	"bytes"
	"io"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-scriptlet/render"
)

//
// RenderWriter
//

type RenderWriter struct {
	originalWriter io.Writer
	context        *commonjs.Context
	render         render.RenderFunc
	buffer         *bytes.Buffer
}

func NewRenderWriter(originalWriter io.Writer, context *commonjs.Context, renderer string) (*RenderWriter, error) {
	if render_, err := render.GetRenderer(renderer); err == nil {
		// Note: renderer can be nil
		return &RenderWriter{
			originalWriter: originalWriter,
			context:        context,
			render:         render_,
			buffer:         bytes.NewBuffer(nil),
		}, nil
	} else {
		return nil, err
	}
}

// ([io.Writer] interface)
func (self *RenderWriter) Write(b []byte) (int, error) {
	if self.render == nil {
		// Optimize for empty renderer
		return self.originalWriter.Write(b)
	} else {
		return self.buffer.Write(b)
	}
}

// [io.StringWriter] interface
func (self *RenderWriter) WriteString(s string) (int, error) {
	if self.render == nil {
		// Optimize for empty renderer
		return render.WriteString(self.originalWriter, s)
	} else {
		return self.buffer.WriteString(s)
	}
}

// [io.ByteWriter] interface
func (self *RenderWriter) WriteByte(c byte) error {
	if self.render == nil {
		// Optimize for empty renderer
		_, err := self.originalWriter.Write([]byte{c})
		return err
	} else {
		return self.buffer.WriteByte(c)
	}
}

// ([io.Closer] interface)
func (self *RenderWriter) Close() error {
	if self.render == nil {
		// Optimize for empty renderer
		return nil
	} else {
		return self.render(self.originalWriter, self.buffer.Bytes(), false, self.context)
	}
}

// ([WrappingWriter] interface)
func (self *RenderWriter) GetWrappedWriter() io.Writer {
	return self.originalWriter
}
