package jst

import (
	"io"
	"strings"

	"github.com/dop251/goja"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-ard"
)

func Present(environment *commonjs.Environment, object *goja.Object, writer io.Writer, variables ard.StringMap) error {
	context := NewContext(writer, variables)
	if _, err := environment.GetAndCall(object, "present", context); err == nil {
		return context.Flush()
	} else {
		return err
	}
}

func RequireAndPresent(environment *commonjs.Environment, id string, writer io.Writer, variables ard.StringMap) error {
	if object, err := environment.Require(id, false, nil); err == nil {
		return Present(environment, object, writer, variables)
	} else {
		return err
	}
}

func RequireAndPresentString(environment *commonjs.Environment, id string, variables ard.StringMap) (string, error) {
	var builder strings.Builder
	if err := RequireAndPresent(environment, id, &builder, variables); err == nil {
		return builder.String(), nil
	} else {
		return "", err
	}
}
