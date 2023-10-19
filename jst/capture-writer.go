package jst

import (
	"bytes"
	"io"
)

type CapturedFunc func(name string, value string)

//
// CaptureWriter
//

type CaptureWriter struct {
	originalWriter io.Writer
	name           string
	onCaptured     CapturedFunc
	buffer         *bytes.Buffer
}

func NewCaptureWriter(originalWriter io.Writer, name string, onCaptured CapturedFunc) *CaptureWriter {
	return &CaptureWriter{
		originalWriter: originalWriter,
		name:           name,
		onCaptured:     onCaptured,
		buffer:         bytes.NewBuffer(nil),
	}
}

// ([io.Writer] interface)
func (self *CaptureWriter) Write(b []byte) (int, error) {
	return self.buffer.Write(b)
}

// [io.StringWriter] interface
func (self *CaptureWriter) WriteString(s string) (int, error) {
	return self.buffer.WriteString(s)
}

// [io.ByteWriter] interface
func (self *CaptureWriter) WriteByte(c byte) error {
	return self.buffer.WriteByte(c)
}

// ([io.Closer] interface)
func (self *CaptureWriter) Close() error {
	self.onCaptured(self.name, self.buffer.String())
	return nil
}

// ([WrappingWriter] interface)
func (self *CaptureWriter) GetWrappedWriter() io.Writer {
	return self.originalWriter
}
