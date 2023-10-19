package jst

import (
	"fmt"
	"io"
	"strings"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-scriptlet/render"
)

// ([scriptlet.RenderFunc] signature)
func RenderJST(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	if !js {
		return nil
	}

	content_, err := render.ToString(content)
	if err != nil {
		return err
	}

	tags, finalChunkStart, err := getTags(content_)
	if err != nil {
		return err
	}

	scriptletContext := NewScriptletContext(writer)

	if err = render.WritePresentStart(writer); err != nil {
		return err
	}

	if len(tags) == 0 {
		// Optimize for when there are no tags
		if err = scriptletContext.AsContextWrite(content_); err != nil {
			return err
		}
	} else {
		last := 0

		for _, tag := range tags {
			// Previous chunk
			if err = scriptletContext.AsContextWrite(content_[last:tag.start]); err != nil {
				return err
			}
			last = tag.end

			code := content_[tag.start+2 : tag.end-2]

			// Swallow trailing newline by default
			swallowTrailingNewline := true

			if content_[tag.end-3] == '/' {
				// Disable the swallowing of trailing newline
				code = code[:len(code)-1]
				swallowTrailingNewline = false
			}

			trimmedCode := strings.TrimSpace(code)
			if trimmedCode == "" {
				// Optimize for empty scriptlets
				continue
			}

			// Handle sugar
			handledSugar := false
			for prefix, handleSugar := range sugarHandlers {
				if strings.HasPrefix(code, prefix) {
					if allowTrailingNewline, err := handleSugar(scriptletContext, prefix, code); err == nil {
						if allowTrailingNewline {
							swallowTrailingNewline = false
						}
						handledSugar = true
						break
					} else {
						return err
					}
				}
			}

			if !handledSugar {
				// Scriptlet tag (no sugar)
				if _, err = render.WriteString(writer, trimmedCode); err != nil {
					return err
				}
				if _, err = render.WriteString(writer, "\n"); err != nil {
					return err
				}
			}

			if swallowTrailingNewline {
				// Skip trailing newline
				if (tag.end <= finalChunkStart) && (content_[tag.end] == '\n') {
					last++
				}
			}
		}

		if last <= finalChunkStart {
			// Final chunk
			if err = scriptletContext.AsContextWrite(content_[last:]); err != nil {
				return err
			}
		}
	}

	return render.WritePresentEnd(writer)
}

//
// tag
//

type tag struct {
	start int
	end   int
}

func getTags(content string) ([]tag, int, error) {
	var tags []tag
	finalChunkStart := len(content) - 1
	start := -1
	var skipNext bool

	for index, rune_ := range content {
		if skipNext {
			skipNext = false
			continue
		}

		switch rune_ {
		case '<':
			// Opening delimiter?
			if (index < finalChunkStart) && (content[index+1] == '%') {
				// Not escaped?
				if (index == 0) || (content[index-1] != '\\') {
					start = index
					skipNext = true
				}
			}

		case '%':
			// Closing delimiter?
			if (index < finalChunkStart) && (content[index+1] == '>') {
				// Not escaped?
				if (index == 0) || (content[index-1] != '\\') {
					skipNext = true
					if start != -1 {
						tags = append(tags, tag{start, index + 2})
						start = -1
					} else {
						return nil, -1, fmt.Errorf("closing delimiter without an opening delimiter at position %d", index)
					}
				}
			}
		}
	}

	return tags, finalChunkStart, nil
}
