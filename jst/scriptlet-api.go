package jst

import (
	contextpkg "context"
	"io"
	"time"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/go-scriptlet/render"
	"github.com/tliron/kutil/util"
)

// ([commonjs.CreateExtensionFunc] signature)
func CreateScriptletExtension(jsContext *commonjs.Context) any {
	return NewScriptletAPI(jsContext)
}

//
// ScriptletAPI
//

type ScriptletAPI struct {
	jsContext *commonjs.Context
}

func NewScriptletAPI(jsContext *commonjs.Context) *ScriptletAPI {
	return &ScriptletAPI{
		jsContext: jsContext,
	}
}

func (self *ScriptletAPI) Render(writer io.Writer, content any, renderer string, timeoutSeconds float64) error {
	return render.Render(writer, content, renderer, false, self.jsContext)
}

func (self *ScriptletAPI) RenderFrom(writer io.Writer, id string, renderer string, timeoutSeconds float64) error {
	context := contextpkg.Background()
	if timeoutSeconds > 0.0 {
		var cancelContext contextpkg.CancelFunc
		context, cancelContext = contextpkg.WithTimeout(context, time.Duration(timeoutSeconds*float64(time.Second)))
		defer cancelContext()
	}

	if url, err := self.jsContext.ResolveAndWatch(context, id, false); err == nil {
		if reader, err := url.Open(context); err == nil {
			reader = util.NewContextualReadCloser(context, reader)
			if err := render.Render(writer, reader, renderer, false, self.jsContext); err == nil {
				return reader.Close()
			} else {
				commonlog.CallAndLogWarning(reader.Close, "ScriptletAPI.RenderFrom", self.jsContext.Environment.Log)
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}

func (self *ScriptletAPI) RenderToString(content any, renderer string) (string, error) {
	return render.RenderToString(content, renderer, false, self.jsContext)
}
