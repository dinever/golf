package Golf

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

var errorTemplate = `<!DOCTYPE HTML>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
    <title>Error: {{.Code}} {{.Title}}</title>
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

    </style>

  </head>
  <body style="background-color: #E4F9F5">
    <div id="header">
      <div class="container">

        <div id="title">
          <h1>Error: {{.Code}} {{.Title}}</h1>
        </div>
      </div>
    </div>
    <div class="container">
      <p>Sorry, the requested URL {{.URL}} caused an error: </p>
      <pre><code>{{.Message}}</code></pre>
      <pre><code>{{.StackTrace}}</code></pre>
    </div>
  </body>
</html>`

var tmpl = template.New("error")

type templateError struct {
	Format     string
	Parameters []interface{}
}

func (e *templateError) Error() string {
	return fmt.Sprintf(e.Format, e.Parameters...)
}

func Errf(format string, parameters ...interface{}) error {
	return &templateError{
		Format:     format,
		Parameters: parameters,
	}
}

// The default error handler
func defaultErrorHandler(ctx *Context) {
	tmpl.Parse(errorTemplate)
	var buf bytes.Buffer
	tmpl.Execute(&buf, map[string]interface{}{
		"Code":    ctx.StatusCode,
		"Title":   http.StatusText(ctx.StatusCode),
		"Message": http.StatusText(ctx.StatusCode),
	})
	ctx.Write(buf.String())
	ctx.Send()
}
