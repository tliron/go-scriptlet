package jst

import (
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonjs-goja/api"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
)

func NewDefaultEnvironment(log commonlog.Logger, urlContext *exturl.Context, basePaths ...exturl.URL) *commonjs.Environment {
	environment := commonjs.NewEnvironment(urlContext, basePaths...)
	if log != nil {
		environment.Log = log
	}
	environment.Precompile = Precompile
	environment.Extensions = CreateDefaultExtensions(true)

	return environment
}

func CreateDefaultExtensions(lateBind bool) []commonjs.Extension {
	return append(
		api.DefaultExtensions{LateBind: lateBind}.Create(),
		commonjs.Extension{
			Name:   "scriptlet",
			Create: CreateScriptletExtension,
		},
	)
}
