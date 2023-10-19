Scriptlets for Go
=================

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/go-scriptlet.svg)](https://github.com/tliron/go-scriptlet/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/tliron/go-scriptlet.svg)](https://pkg.go.dev/github.com/tliron/go-scriptlet)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/go-scriptlet)](https://goreportcard.com/report/github.com/tliron/go-scriptlet)

A 100% Go extensible library for rendering text from JavaScript Templates (JST) by executing embedded
scriptlets. Supports "sugar" via scriptlet codes for rendering expressions, capturing blocks, importing
files, and more. Extend with your own sugar.

Useful for applications that generate text dynamically, such as web pages, documentation, and configuration
files.

For example, render this:

    <div>
    <% for (let i = 0; i < 3; i++) { %>
        <div>Number <%= i+1 %></div>
    <% } %>
    </div>

to this:

    <div>
        <div>Number 1</div>
        <div>Number 2</div>
        <div>Number 3</div>
    </div>

For a comprehensive web platform built with go-scriptlet, see
[Prudence](https://github.com/tliron/prudence).

JavaScript is run in a [CommonJS-style](https://wiki.commonjs.org/wiki/CommonJS) modular environment
via the [Goja](https://github.com/dop251/goja) JavaScript engine (100% Go). See
[commonjs-goja](https://github.com/tliron/commonjs-goja) for the full implementation.

The rendering API can actually support other engines, not just JST. Included is support for rendering
Markdown (via [goldmark](https://github.com/yuin/goldmark)), HTML sanitizing (via
[bluemonday](https://github.com/microcosm-cc/bluemonday)), as well as minifying various web formats:
HTML, CSS, JSON, XML, web JavaScript, and SVG.

Basic Usage
===========

This is the minimum code necessary to render JST templates:

```go
import (
	"fmt"
	"os"
	"github.com/tliron/exturl"
	"github.com/tliron/go-scriptlet/jst"
	"github.com/tliron/go-scriptlet/markdown"
	"github.com/tliron/go-scriptlet/minify"
	"github.com/tliron/go-scriptlet/sanitize"
)

func init() {
	jst.RegisterDefaultRenderers()
	jst.RegisterDefaultSugar()
	sanitize.RegisterDefaultRenderers()
	markdown.RegisterDefaultRenderers()
	minify.RegisterDefaultRenderers()
}

func main() {
	path := os.Args[1] // the path of the template to render

	urlContext := exturl.NewContext()
	defer urlContext.Release()

	wd, _ := urlContext.NewWorkingDirFileURL() // our base path is the working dir

	environment := jst.NewDefaultEnvironment(nil, urlContext, wd)
	defer environment.Release()

	if err := jst.Present(environment, id, os.Stdout, nil); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
```

Included is an example of a [very simple web server](examples/web.go) serving JST dynamically.
To see it in action, clone this repository and run:

    go run examples/web.go

Direct your browser to [http://localhost:8080](http://localhost:8080) to see the examples.

A real-world live scriptlet-based web server would probably want to cache some content instead
of rendering it dynamically for each request. That feature is beyond the scope of go-scriptlet.
Again, check out [Prudence](https://github.com/tliron/prudence) for a comprehensive platform.

JavaScript Templates (JST)
==========================

The scriptlet delimiters are `<%` and `%>`. Characters that follow right after the opening
delimiter specify special "sugar".

The delimiters can be escaped by prefixing a backslash: `\<%` and `\%>`.

Note that scriptlets are not self-contained programs, and indeed allow you to mix JavaScript
with literal text:

    <% for (let i = 0; i < 10; i++) { %>
        <p>Number <%= i %></p>
    <% } %>

By default, if a scriptlet's end delimiter also ends a line then the renderer will "swallow"
the trailing newline. This helps you avoid cluttering your output with empty lines, and is
quite intuitive (for some people, at least). For example, this template:

    No
    <% var x = 1; %>
    empty
    <% x += 1; %>
    lines!

will be rendered as this:

    No
    Empty
    Lines

To disable this feature use `/%>` as the closing delimiter. This:

    Empty
    <% var x = 1; /%>
    line!

will be output as this:

    Empty

    line!

The one exception is the "expression" sugar, `<%=`, which does *not* swallow the trailing newline
because it's intended to be used within flows of text.

Variables
---------

When you declare a JavaScript variable with `var`, `let`, or `const` in JST it is not a true
JavaScript global, because the implementation wraps the entire JST script in a function. If you
do need variables to cross the boundary of a single JST file, go-scriptlet does come with two
solutions.

First, there context-local variables, which are also exposed in Go, that would persist even if you
compose a more complex render pipeline around a single context, e.g. using the "embed" sugar
(see below), HTTP middleware injecting or using a variable, etc. To access them in JavaScript:

```javascript
this.variables.myVar = 'local';
```

The underlying CommonJS library also provides true globals that persist across *all* contexts:

```javascript
env.variables.myVar = 'global';
```

Note that access to this namespace is thread-safe, but individual variables are not.

Built-in Sugar
--------------

### Comment: `<%# anything %>`

The content of the scriptlet is ignored. Can be useful for quickly disabling other scriptlets
during development. Example:

    <%#
    This is
    ignored (and also not rendered into the script)
    %>

Note that this does not even insert JavaScript comments, so though it's functionally equivalent
to JavaScript comments, it is not identical in implementation:

    <%
    // This is
    // ignored (but still rendered)
    %>

### Expression: `<%= expr %>` or `<%=/ expr %>`

Write a JavaScript expression. Example:

    Hello! Your name is <%= ['linus', 'torvalds'].join(' ').toUpperCase() %>!

The "/" variant will escape HTML characters before writing.

### Variable: `<%== string_expr,... %>` or `<%==/ string_expr,... %>`

Write a context-local variable. Example:

    <% this.variables.name = 'Linus'; %>
    Hello! Your name is <%== 'name' %>!

Also supports safely accessing nested variables by providing an array of strings:

    <% this.variables.person = {name: 'Linus'}; %>
    Hello! Your name is <%== 'person', 'name' %>!

The "/" variant will escape HTML characters before writing.

### Insert: `<%+ string_expr %>` or `<%+ string_expr, string_expr %>`

Loads and writes the contents of a file, optionally rendering it. Simple text insert:

    <%+ '../docs/README.md' %>

With rendering:

    <%+ '../docs/README.md', 'markdown' %>

Complete URLs are supported:

    <%+ 'https://raw.githubusercontent.com/tliron/go-scriptlet/main/README.md', 'markdown' %>

The short form, without rendering, is optimized to not load the entire file
into memory, instead doing a buffered copy to the output stream.

### Embed: `<%& string_expr %>`

Renders a JST file. The embedded JST gets a copy of all the parent's context-local variables,
but changes to the variables are not reflected back to the parent in order to ensure data
consistency. Example:

    <% this.variables.name = 'Linus'; %>
    <%& './header.jst' %>

Where `header.jst` can be this:

   Your name is <%== 'name' %>

### Capture: `<%! string_expr %>` and `<%!!%>`

Captures the enclosed text into a context-local variable. Does *not* write it. Example:

    <%! 'greeting' %>
    <div>
        Hello, <%==/ 'name' %>!
    </div>
    <%!!%>

    The greeting is: <%==/ 'greeting' %>

When used in conjunction with "embed" sugar you can make page templates:

    <%! 'body' %>
    Hello, <%==/ 'name' %>!
    <%!!%>
    <%& './page.jst' %>

Where `page.jst` can be this:

    <html>
    <body>
        <%==/ 'body' %>
    </body>

### Render: `<%^ string_expr %>` and `<%^^%>`

Renders the enclosed text *before* writing it. Example:

    <%^ 'markdown' %>
    This is Markdown
    ================

    Hello, <%== 'name' %>!

    It is a *markup* language for generating HTML.
    <%^^%>

Note that, as in this example, any other JST scriptlets inside the enclosed text are executed
as usual.

Default renderers:

* `sanitizehtml`
* `markdown` or `md`
* `extendedmarkdown` or `extendedmd`
* `mincss`
* `minhtml`
* `minsvg`
* `minjs`
* `minjson`
* `minxml`

Custom Sugar
------------

If the built-in sugar is not sweet enough for you then you can add your own.

Your custom sugar is registered on a prefix, which is a string that will be checked against what
immediately follows the `<%` opening delimiter. Note that not only must it be unique so that it
won't overlap with other sugar, but also that it should be unambiguous. Thus you shouldn't register
both the the `-` and the `->` prefixes because the former is included in the latter.

Your sugar implementation has three arguments, a `ScriptletContext`, the prefix, and the raw text
between the two scriptlet delimiters (which includes your prefix). Your implementation can do anything,
but what it most likely will do is write JavaScript source code into the context. Included are utility
functions to help you do this.

Example:

```go
func init() {
	jst.RegisterSugar("~", HandleInBed)
}

// ([jst.HandleSugarFunc] signature)
func HandleInBed(scriptletContext *jst.ScriptletContext, prefix string, code string) (bool, error) {
	code = code[len(prefix):] // skip the "~" prefix
	code = strings.TrimSpace(code) // remove spaces on each side
	code += " in bed" // sweet, sweet sugar
	return false, scriptletContext.AsContextWrite(code) // this.write('...');
}
```

And then using it in JST:

    <div>
        <%~ I like to watch TV %>
    </div>

Custom Renderers
----------------

The go-scriptlet renderer API is quite straightforward: it accepts a "content" input and writes
to a `io.Writer`. What the [renderer](../render/README.md) actually does, of course, can be quite
sophisticated, as in the case of JST. It could be an entire language implementation. Note that
data streaming is supported, too, because "content" can be a `io.Reader`. Utilities are provided
to help work with inputs and outputs of various types.

Here's a trivial example:

```go
import (
	"io"
	"strings"
	"github.com/tliron/commonjs-goja"
	"github.com/tliron/go-scriptlet/render"
)

func init() {
	render.RegisterRenderer("doublespace", RenderDoubleSpace)
}

// ([render.RenderFunc] signature)
func RenderDoubleSpace(writer io.Writer, content any, js bool, jsContext *commonjs.Context) error {
	if content_, err := render.ToString(content); err == nil {
		content_ = strings.ReplaceAll(content_, " ", "  ")
		if js {
			return render.AsPresenter(writer, content_)
		} else {
			_, err = render.WriteString(writer, content_)
			return err
		}
	} else {
		return err
	}
}
```
