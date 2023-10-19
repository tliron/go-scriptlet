package scriptlet_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/exturl"
	"github.com/tliron/go-scriptlet/jst"
	"github.com/tliron/go-scriptlet/jst/sugar"
	"github.com/tliron/go-scriptlet/markdown"
	"github.com/tliron/go-scriptlet/minify"
	"github.com/tliron/go-scriptlet/sanitize"

	_ "github.com/tliron/commonlog/simple"
)

func init() {
	jst.DebugWriter = os.Stdout

	jst.RegisterDefaultRenderers()
	sanitize.RegisterDefaultRenderers()
	markdown.RegisterDefaultRenderers()
	minify.RegisterDefaultRenderers()

	sugar.RegisterDefaultSugar()
	jst.RegisterSugar("~", handleInBed)
}

func TestScriptlet(t *testing.T) {
	urlContext := exturl.NewContext()
	defer urlContext.Release()

	path := filepath.Join(getRoot(t), "examples")

	environment := jst.NewDefaultEnvironment(nil, urlContext, urlContext.NewFileURL(path))
	defer environment.Release()

	environment.Precompile = testPrecompile

	testPresent(t, environment, "./broken.jst", false)
	testPresent(t, environment, "./no-scriptlets.jst", true)

	testPresent(t, environment, "./module.jst", true)
	testPresent(t, environment, "./newline.jst", true)
	testPresent(t, environment, "./variables.jst", true)
	testPresent(t, environment, "./comments.jst", true)
	testPresent(t, environment, "./console.jst", true)
	testPresent(t, environment, "./capture.jst", true)
	testPresent(t, environment, "./render.jst", true)
	testPresent(t, environment, "./insert.jst", true)
	testPresent(t, environment, "./embed.jst", true)
	testPresent(t, environment, "./custom.jst", true)

	testPresent(t, environment, "./example.md", true)
	testPresent(t, environment, "./example.mincss", true)
}

func testPresent(t *testing.T, environment *commonjs.Environment, id string, success bool) {
	if content, err := jst.RequireAndPresentString(environment, id, nil); err == nil {
		printTitle("Render for " + id)
		printContent(content)
	} else if success {
		t.Errorf("engine.Present: %s", err)
	}
}

// ([jst.HandleSugarFunc] signature)
func handleInBed(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):]                           // skip the "~" prefix
	code = strings.TrimSpace(code)                      // remove spaces on each side
	code += " in bed"                                   // sweet, sweet sugar
	return false, scriptletContext.AsContextWrite(code) // this.write('...');
}

func getRoot(t *testing.T) string {
	var root string
	var ok bool
	if root, ok = os.LookupEnv("SCRIPTLET_TEST_ROOT"); !ok {
		var err error
		if root, err = os.Getwd(); err != nil {
			t.Errorf("os.Getwd: %s", err.Error())
		}
	}
	return root
}

// ([commonjs.PrecompileFunc] signature)
func testPrecompile(url exturl.URL, script string, context *commonjs.Context) (string, error) {
	printTitle("Script for " + url.String())
	if script, err := jst.Precompile(url, script, context); err == nil {
		return script, nil
	} else {
		return "", err
	}
}

func printTitle(title string) {
	fmt.Println()
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", len(title)))
}

func printContent(content string) {
	fmt.Println(strings.TrimRight(content, "\n"))
}
