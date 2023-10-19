package render

import (
	"io"
)

func WritePresentStart(writer io.Writer) error {
	_, err := WriteString(writer, "exports.present = function() {\n")
	return err
}

func WritePresentEnd(writer io.Writer) error {
	_, err := WriteString(writer, "};\n")
	return err
}

func AsPresenter(writer io.Writer, content any) error {
	if err := WritePresentStart(writer); err != nil {
		return err
	}
	if err := AsContextWrite(writer, content); err != nil {
		return err
	}
	return WritePresentEnd(writer)
}

func AsContextWrite(writer io.Writer, content any) error {
	if content == nil {
		return nil
	}

	if _, err := WriteString(writer, "this.write('"); err != nil {
		return err
	}

	runeReader := ToRuneReader(content)
	for {
		rune_, _, err := runeReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch rune_ {
		case '\n':
			if _, err := WriteString(writer, "\\n"); err != nil {
				return err
			}

		case '\'':
			if _, err := WriteString(writer, "\\'"); err != nil {
				return err
			}

		case '\\':
			if _, err := WriteString(writer, "\\\\"); err != nil {
				return err
			}

		default:
			if _, err := WriteRune(writer, rune_); err != nil {
				return err
			}
		}
	}

	_, err := WriteString(writer, "');\n")
	return err
}
