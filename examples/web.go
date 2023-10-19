package main

import (
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tliron/commonjs-goja"
	"github.com/tliron/commonlog"
	"github.com/tliron/exturl"
	"github.com/tliron/go-ard"
	"github.com/tliron/go-scriptlet/jst"
	"github.com/tliron/go-scriptlet/jst/sugar"
	"github.com/tliron/go-scriptlet/markdown"
	"github.com/tliron/go-scriptlet/minify"
	"github.com/tliron/go-scriptlet/sanitize"
	"github.com/tliron/kutil/util"

	_ "github.com/tliron/commonlog/simple"
)

var address = ":8080"
var dir string
var environment *commonjs.Environment
var log = commonlog.GetLogger("web")

func init() {
	jst.RegisterDefaultRenderers()
	sanitize.RegisterDefaultRenderers()
	markdown.RegisterDefaultRenderers()
	minify.RegisterDefaultRenderers()

	sugar.RegisterDefaultSugar()
	jst.RegisterSugar("~", handleInBed)
}

func main() {
	util.ExitOnSignals()
	util.InitializeColorization("true")
	commonlog.Configure(2, nil)

	_, dir, _, _ = runtime.Caller(0)
	dir = filepath.Dir(dir)

	log.Noticef("serving directory: %s", dir)

	urlContext := exturl.NewContext()
	util.OnExitError(urlContext.Release)

	environment = jst.NewDefaultEnvironment(nil, urlContext, urlContext.NewFileURL(dir))
	util.OnExitError(environment.Release)

	environment.OnFileModified = func(id string, module *commonjs.Module) {
		if module != nil {
			environment.Log.Infof("module changed: %s", module.Id)
		} else if id != "" {
			environment.Log.Infof("file changed: %s", id)
		}

		// Note: this clears the *program* cache for commonjs-goja,
		// *not* a rendered HTTP content cache (which we don't have)
		environment.ClearCache()
	}

	util.FailOnError(environment.StartWatcher())

	http.HandleFunc("/", handleHttp)
	log.Noticef("starting server: %s", address)
	util.FailOnError(http.ListenAndServe(address, nil))
}

func handleHttp(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	response.Header().Add("Content-Type", "text/html; charset=UTF-8")
	response.Header().Add("Content-Language", "en")

	// TODO: handle io.WriteString errors
	io.WriteString(response, "<!DOCTYPE html>\n")
	io.WriteString(response, "<html lang=\"en\">\n")
	io.WriteString(response, "<head><meta charset=\"UTF-8\" /><link rel=\"icon\" href=\"data:,\" /><title>go-scriptlet</title></head>\n")
	io.WriteString(response, "<body>\n")
	io.WriteString(response, "<h1>"+html.EscapeString(path)+"</h1>\n")

	if err := present(response, request.URL.Path); err != nil {
		environment.Log.Errorf("%s", err.Error())
		io.WriteString(response, "<p style=\"color: red\">"+html.EscapeString(err.Error())+"</p>\n")
	}

	io.WriteString(response, "</body\n")
	io.WriteString(response, "</html>\n")
}

func present(writer io.Writer, path string) error {
	if path == "/" {
		// Directory
		log.Infof("directory: %s", path)
		if dirEntries, err := os.ReadDir(dir); err == nil {
			for _, dirEntry := range dirEntries {
				if !dirEntry.IsDir() {
					name := dirEntry.Name()
					switch filepath.Ext(name) {
					case ".go", ".js":
						// Skip
					default:
						name = html.EscapeString(name)
						if _, err := io.WriteString(writer, "<a href=\""+name+"\">"+name+"</a><br/>\n"); err != nil {
							return err
						}
					}
				}
			}

			return nil
		} else {
			return err
		}
	} else {
		log.Infof("presenting: %s", path)
		// Set "path" in this.variables
		return jst.RequireAndPresent(environment, filepath.Join(".", path), writer, ard.StringMap{"path": path})
	}
}

// ([jst.HandleSugarFunc] signature)
func handleInBed(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):]                           // skip the "~" prefix
	code = strings.TrimSpace(code)                      // remove spaces on each side
	code += " in bed"                                   // sweet, sweet sugar
	return false, scriptletContext.AsContextWrite(code) // this.write('...');
}
