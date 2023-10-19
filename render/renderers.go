package render

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/tliron/commonjs-goja"
)

// Special handling for string, []byte, and [io.Reader].
type RenderFunc func(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error

var renderers = make(map[string]RenderFunc)

// The render function can be nil, in which case no rendering will occur.
func RegisterRenderer(renderer string, render RenderFunc) {
	renderers[renderer] = render
}

func GetRenderers() []string {
	names := make([]string, len(renderers))
	index := 0
	for name := range renderers {
		names[index] = name
		index++
	}
	return names
}

func GetRenderer(renderer string) (RenderFunc, error) {
	if renderer == "" {
		// Empty string means nil renderer
		return nil, nil
	} else if render, ok := renderers[renderer]; ok {
		// Note that the renderer can be nil
		return render, nil
	} else {
		return nil, fmt.Errorf("unsupported renderer: %s", renderer)
	}
}

// Special handling for string, []byte, and [io.Reader].
func Render(writer io.Writer, content any, renderer string, js bool, jsContext *commonjs.Context) error {
	if render, err := GetRenderer(renderer); err == nil {
		if render == nil {
			if js {
				return AsPresenter(writer, content)
			} else {
				return nil
			}
		} else {
			return render(writer, content, js, jsContext)
		}
	} else {
		return err
	}
}

// Special handling for string, []byte, and [io.Reader].
func RenderToBytes(content any, renderer string, js bool, jsContext *commonjs.Context) ([]byte, error) {
	var buffer bytes.Buffer
	if err := Render(&buffer, content, renderer, js, jsContext); err == nil {
		return buffer.Bytes(), nil
	} else {
		return nil, err
	}
}

// Special handling for string, []byte, and [io.Reader].
func RenderToString(content any, renderer string, js bool, jsContext *commonjs.Context) (string, error) {
	var builder strings.Builder
	if err := Render(&builder, content, renderer, js, jsContext); err == nil {
		return builder.String(), nil
	} else {
		return "", err
	}
}
