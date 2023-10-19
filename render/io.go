package render

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/tliron/kutil/util"
)

// Special handling for [io.Reader].
func ToString(content any) (string, error) {
	switch content_ := content.(type) {
	case io.Reader:
		var builder strings.Builder
		if _, err := io.Copy(&builder, content_); err == nil {
			return builder.String(), nil
		} else {
			return "", err
		}

	default:
		return util.ToString(content_), nil
	}
}

// Special handling for [io.Reader].
func ToBytes(content any) ([]byte, error) {
	switch content_ := content.(type) {
	case []byte:
		return content_, nil

	case io.Reader:
		return io.ReadAll(content_)

	default:
		return util.StringToBytes(util.ToString(content_)), nil
	}
}

// Special handling for nil, string, []byte, and [io.Reader].
func Write(writer io.Writer, content any) (int, error) {
	switch content_ := content.(type) {
	case nil:
		return 0, nil

	case string:
		return WriteString(writer, content_)

	case []byte:
		return writer.Write(content_)

	case io.Reader:
		n, err := io.Copy(writer, content_)
		return int(n), err

	default:
		return WriteString(writer, util.ToString(content_))
	}
}

func WriteString(writer io.Writer, content string) (int, error) {
	return io.WriteString(writer, content)

	/*
		if stringWriter, ok := writer.(io.StringWriter); ok {
			return stringWriter.WriteString(s)
		}
		// Note: this will break if the writer is async, because the
		// underlying bytes might change before the string is accessed
		return writer.Write(util.StringToBytes(s))
	*/
}

func WriteRune(writer io.Writer, rune_ rune) (int, error) {
	return WriteString(writer, string(rune_))
}

func ToReader(content any) io.Reader {
	switch content_ := content.(type) {
	case io.Reader:
		return content_

	case []byte:
		return bytes.NewReader(content_)

	case string:
		return strings.NewReader(content_)

	default:
		return strings.NewReader(util.ToString(content_))
	}
}

func ToRuneReader(content any) io.RuneReader {
	switch content_ := content.(type) {
	case io.Reader:
		return bufio.NewReader(content_)

	case []byte:
		return bytes.NewReader(content_)

	case string:
		return strings.NewReader(content_)

	default:
		return strings.NewReader(util.ToString(content_))
	}
}
