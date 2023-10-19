package jst

type HandleSugarFunc func(scriptletContext *ScriptletContext, prefix string, code string) (bool, error) // return true to allow trailing newlines

var sugarHandlers = make(map[string]HandleSugarFunc)

func RegisterSugar(prefix string, handle HandleSugarFunc) {
	sugarHandlers[prefix] = handle
}
