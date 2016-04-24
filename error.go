package golf

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"reflect"
	"runtime"
	"strconv"
)

const errorTemplate = `<!DOCTYPE HTML><html><head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>Error: {{ .Code }} {{ .Title }}</title>
    <style type="text/css" media="screen">
html,body{
  padding:0;
  margin:0;
  font-family: Tahoma;
  color: #34495e;
}
h1 {
  color: #fff;
  margin: 0;
}
.container {
  max-width: 1220px;
  margin: 0 auto;
  padding: 0 20px;
}
#header {
  display: block;
  background-color: #3498db;
  height: 120px;
  width: 100%;
}
#title {
  padding: 40px 0;
}
.error {
  color: #c0392b;
}
pre code {
  font-family: "Lucida Console", Monaco, monospace;
}
pre.request-dump {
  background-color: #eeeeee;
  padding: 20px 20px 0 20px;
  overflow: auto;
}
#backtrace {
  list-style: none;
  padding-left: 0;
  font-family: "Lucida Console", Monaco, monospace;
}
#backtrace li {
  border-left: 5px solid #61A8DC;
  padding-left: 20px;
}
#backtrace .file {
  color: #61A8DC;
}
#backtrace .lineno {
color: #ff8a00;
}
#backtrace .method {
  color: #34a853;
}

    </style>
    <body>
      <div id="header">
        <div class="container">
          <div id="title">
            <h1>Error: {{ .Code }} {{ .Title }}</h1>
          </div>
        </div>
      </div>
      <div class="container">
        <p>Sorry, the requested URL  caused an error: </p>
        <pre><code>{{ .Message }}</code></pre>
        <h2>HTTP Request</h2>
        <pre class="request-dump">{{ .HTTPRequest }}</pre>
        {{ if .StackTrace }}
        <h2>Traceback</h2>
        <ul id="backtrace">
          {{ range .StackTrace }}
          <li><p><span class="file">{{ .File }}:</span><span class="lineno">{{ .Number }}</span> <span class="method">{{ .Method }}</span></p></li>
          {{ end }}
        </ul>
        {{ end }}
      </div>

    </body>
</html>`

const maxFrames = 20

var tmpl = template.New("error")

// The default error handler
func defaultErrorHandler(ctx *Context, data ...map[string]interface{}) {
	var renderData map[string]interface{}
	if len(data) == 0 {
		renderData = make(map[string]interface{})
		renderData["Code"] = ctx.statusCode
		renderData["Title"] = http.StatusText(ctx.statusCode)
		renderData["Message"] = http.StatusText(ctx.statusCode)
	} else {
		renderData = data[0]
	}
	if _, ok := renderData["Code"]; !ok {
		renderData["Code"] = ctx.statusCode
	}
	if _, ok := renderData["Title"]; !ok {
		renderData["Title"] = http.StatusText(ctx.statusCode)
	}
	if _, ok := renderData["Message"]; !ok {
		renderData["Message"] = http.StatusText(ctx.statusCode)
	}
	httpRequest, _ := httputil.DumpRequest(ctx.Request, true)
	renderData["HTTPRequest"] = string(httpRequest)
	var buf bytes.Buffer
	tmpl.Parse(errorTemplate)
	tmpl.Execute(&buf, renderData)
	ctx.Send(&buf)
}

// Frame represent a stack frame inside of a Honeybadger backtrace.
type Frame struct {
	Number string `json:"number"`
	File   string `json:"file"`
	Method string `json:"method"`
}

// Error provides more structured information about a Go error.
type Error struct {
	err     interface{}
	Message string
	Class   string
	Stack   []*Frame
}

// Error returns the error message
func (e Error) Error() string {
	return e.Message
}

// StackTraceString returns the stack trace in a string format.
func (e Error) StackTraceString() string {
	buf := new(bytes.Buffer)
	for _, v := range e.Stack {
		fmt.Fprintf(buf, "%s: %s\n\t%s\n", v.File, v.Number, v.Method)
	}
	return string(buf.Bytes())
}

// Errorf returns an templateError.
func Errorf(format string, parameters ...interface{}) error {
	return fmt.Errorf(format, parameters...)
}

// NewError creates a new error instance
func NewError(msg interface{}) Error {
	var err error

	switch t := msg.(type) {
	case Error:
		return t
	case error:
		err = t
	default:
		err = fmt.Errorf("%v", t)
	}

	return Error{
		err:     err,
		Message: err.Error(),
		Class:   reflect.TypeOf(err).String(),
		Stack:   generateStack(3),
	}
}

func generateStack(offset int) (frames []*Frame) {
	stack := make([]uintptr, maxFrames)
	length := runtime.Callers(2+offset, stack[:])
	for _, pc := range stack[:length] {
		f := runtime.FuncForPC(pc)
		if f == nil {
			continue
		}
		file, line := f.FileLine(pc)
		frame := &Frame{
			File:   file,
			Number: strconv.Itoa(line),
			Method: f.Name(),
		}
		frames = append(frames, frame)
	}

	return
}
