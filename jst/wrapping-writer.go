package jst

import (
	"io"
)

//
// WrappingWriter
//

type WrappingWriter interface {
	io.WriteCloser

	GetWrappedWriter() io.Writer
}

func UnwrapWriters(writer io.Writer) (io.Writer, error) {
	for true {
		if wrappingWriter, ok := writer.(WrappingWriter); ok {
			if err := wrappingWriter.Close(); err == nil {
				writer = wrappingWriter.GetWrappedWriter()
			} else {
				return writer, err
			}
		} else {
			break
		}
	}

	return writer, nil
}
