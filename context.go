package Golf

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Context is a wrapper of http.Request and http.ResponseWriter.
type Context struct {
	// http.Request
	Request *http.Request

	// http.ResponseWriter
	Response http.ResponseWriter

	// URL Parameter
	Params map[string]string

	// HTTP status code
	StatusCode int

	// HTTP header as a map
	Header map[string]string

	// HTTP response body as a byte string
	Body []byte

	// The application
	App *Application

	// Data used for sharing values between middlewares
	Data map[string]interface{}

	// Indicating if the response is already sent.
	IsSent bool

	// Indicating loader of the template
	templateLoader string
}

// NewContext creates a Golf.Context instance.
func NewContext(req *http.Request, res http.ResponseWriter, app *Application) *Context {
	ctx := new(Context)
	ctx.Request = req
	ctx.Response = res
	ctx.App = app
	ctx.Header = make(map[string]string)
	ctx.StatusCode = 200
	ctx.Header["Content-Type"] = "text/html;charset=UTF-8"
	ctx.Request.ParseForm()
	ctx.IsSent = false
	ctx.Data = make(map[string]interface{})
	return ctx
}

// Query method retrieves the form data, return empty string if not found.
func (ctx *Context) Query(key string, index ...int) (string, error) {
	if val, ok := ctx.Request.Form[key]; ok {
		if len(index) == 1 {
			return val[index[0]], nil
		}
		return val[0], nil
	}
	return "", errors.New("Query key not found.")
}

// Param method retrieves the parameters from url
// If the url is /:id/, then id can be retrieved by calling `ctx.Param(id)`
func (ctx *Context) Param(key string) (string, error) {
	if val, ok := ctx.Params[key]; ok {
		return val, nil
	}
	return "", errors.New("Parameter not found.")
}

// Redirect method sets the response as a 301 redirection.
// If you need a 302 redirection, please do it by setting the Header manually.
func (ctx *Context) Redirect(url string) {
	ctx.Header["Location"] = url
	ctx.StatusCode = 301
}

// SetCookie set cookies for the request. If expire is 0, create a session cookie.
func (ctx *Context) SetCookie(key string, value string, expire int) {
	now := time.Now()
	cookie := &http.Cookie{
		Name:   key,
		Value:  value,
		Path:   "/",
		MaxAge: expire,
	}
	if expire != 0 {
		expireTime := now.Add(time.Duration(expire) * time.Second)
		cookie.Expires = expireTime
	}
	http.SetCookie(ctx.Response, cookie)
}

// JSON Sends a JSON response.
func (ctx *Context) JSON(obj interface{}) {
	json, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	ctx.Body = json
	ctx.Header["Content-Type"] = "application/json"
}

// Send the response immediately. Set `ctx.IsSent` to `true` to make
// sure that the response won't be sent twice.
func (ctx *Context) Send() {
	if ctx.IsSent {
		return
	}
	for name, value := range ctx.Header {
		ctx.Response.Header().Set(name, value)
	}
	ctx.Response.WriteHeader(ctx.StatusCode)
	ctx.Response.Write(ctx.Body)
	ctx.IsSent = true
}

// Write text on the response body.
func (ctx *Context) Write(content string) {
	ctx.Body = []byte(content)
}

// Abort method returns an HTTP Error by indicating the status code, the corresponding
// handler inside `App.errorHandler` will be called, if user does not set
// the corresponding error handler, the defaultErrorHandler will be called.
func (ctx *Context) Abort(statusCode int) {
	ctx.StatusCode = statusCode
	ctx.App.handleError(ctx, statusCode)
}

// Loader method sets the template loader for this context. This should be done before calling
// `ctx.Render`.
func (ctx *Context) Loader(name string) *Context {
	ctx.templateLoader = name
	return ctx
}

// Render a template file using the built-in Go template engine.
func (ctx *Context) Render(file string, data interface{}) {
	content, e := ctx.App.View.Render(ctx.templateLoader, file, data)
	if e != nil {
		panic(e)
	}
	ctx.Body = []byte(content)
}

// RenderFromString renders a input string.
func (ctx *Context) RenderFromString(tplSrc string, data interface{}) {
	content, e := ctx.App.View.RenderFromString(ctx.templateLoader, tplSrc, data)
	if e != nil {
		panic(e)
	}
	ctx.Body = []byte(content)
}
