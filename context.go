package golf

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Context is a wrapper of http.Request and http.ResponseWriter.
type Context struct {
	// http.Request
	Request *http.Request

	// http.ResponseWriter
	Response http.ResponseWriter

	// URL Parameter
	Params Parameter

	// HTTP status code
	statusCode int

	// The application
	App *Application

	// Session instance for the current context.
	Session Session

	// Indicating if the response is already sent.
	IsSent bool

	// Indicating loader of the template
	templateLoader string
}

// NewContext creates a Golf.Context instance.
func NewContext(req *http.Request, res http.ResponseWriter, app *Application) *Context {
	ctx := new(Context)
	ctx.Request = req
	ctx.Request.ParseForm()
	ctx.Response = res
	ctx.App = app
	ctx.statusCode = 200
	//	ctx.Header["Content-Type"] = "text/html;charset=UTF-8"
	ctx.Request.ParseForm()
	ctx.IsSent = false
	return ctx
}

func (ctx *Context) reset() {
	ctx.statusCode = 200
	ctx.IsSent = false
}

func (ctx *Context) generateSession() Session {
	s, err := ctx.App.SessionManager.NewSession()
	if err != nil {
		return nil
	}
	// Session lifetime should be configurable.
	ctx.SetCookie("sid", s.SessionID(), 3600)
	return s
}

func (ctx *Context) retrieveSession() {
	var s Session
	sid, err := ctx.Cookie("sid")
	if err != nil {
		s = ctx.generateSession()
	} else {
		s, err = ctx.App.SessionManager.Session(sid)
		if err != nil {
			s = ctx.generateSession()
		}
	}
	ctx.Session = s
}

// SendStatus takes an integer and sets the response status to the integer given.
func (ctx *Context) SendStatus(statusCode int) {
	ctx.statusCode = statusCode
	ctx.Response.WriteHeader(statusCode)
}

// StatusCode returns the status code that golf has sent.
func (ctx *Context) StatusCode() int {
	return ctx.statusCode
}

// SetHeader sets the header entries associated with key to the single element value. It replaces any existing values associated with key.
func (ctx *Context) SetHeader(key, value string) {
	ctx.Response.Header().Set(key, value)
}

// AddHeader adds the key, value pair to the header. It appends to any existing values associated with key.
func (ctx *Context) AddHeader(key, value string) {
	ctx.Response.Header().Add(key, value)
}

// Header gets the first value associated with the given key. If there are no values associated with the key, Get returns "".
func (ctx *Context) Header(key string) string {
	return ctx.Request.Header.Get(key)
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
func (ctx *Context) Param(key string) string {
	val, _ := ctx.Params.ByName(key)
	return val
}

// Redirect method sets the response as a 301 redirection.
// If you need a 302 redirection, please do it by setting the Header manually.
func (ctx *Context) Redirect(url string) {
	ctx.SetHeader("Location", url)
	ctx.SendStatus(301)
}

// Cookie returns the value of the cookie by indicating the key.
func (ctx *Context) Cookie(key string) (string, error) {
	c, err := ctx.Request.Cookie(key)
	if err != nil {
		return "", err
	}
	return c.Value, nil
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
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Send(json)
}

// JSONIndent Sends a JSON response, indenting the JSON as desired.
func (ctx *Context) JSONIndent(obj interface{}, prefix, indent string) {
	jsonIndented, err := json.MarshalIndent(obj, prefix, indent)
	if err != nil {
		panic(err)
	}
	ctx.SetHeader("Content-Type", "application/json")
	ctx.Send(jsonIndented)
}

// Send the response immediately. Set `ctx.IsSent` to `true` to make
// sure that the response won't be sent twice.
func (ctx *Context) Send(body interface{}) {
	if ctx.IsSent {
		return
	}
	switch body.(type) {
	case []byte:
		ctx.Response.Write(body.([]byte))
	case string:
		ctx.Response.Write([]byte(body.(string)))
	case *bytes.Buffer:
		ctx.Response.Write(body.(*bytes.Buffer).Bytes())
	default:
		panic(fmt.Errorf("Body type not supported."))
	}
	ctx.IsSent = true
}

func (ctx *Context) requestHeader(key string) string {
	if values, _ := ctx.Request.Header[key]; len(values) > 0 {
		return values[0]
	}
	return ""
}

// ClientIP implements a best effort algorithm to return the real client IP, it parses
// X-Real-IP and X-Forwarded-For in order to work properly with reverse-proxies such us: nginx or haproxy.
// This method is taken from https://github.com/gin-gonic/gin
func (ctx *Context) ClientIP() string {
	clientIP := strings.TrimSpace(ctx.requestHeader("X-Real-Ip"))
	if len(clientIP) > 0 {
		return clientIP
	}
	clientIP = ctx.requestHeader("X-Forwarded-For")
	if index := strings.IndexByte(clientIP, ','); index >= 0 {
		clientIP = clientIP[0:index]
	}
	clientIP = strings.TrimSpace(clientIP)
	if len(clientIP) > 0 {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.Request.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

// Abort method returns an HTTP Error by indicating the status code, the corresponding
// handler inside `App.errorHandler` will be called, if user has not set
// the corresponding error handler, the defaultErrorHandler will be called.
func (ctx *Context) Abort(statusCode int, data ...map[string]interface{}) {
	ctx.App.handleError(ctx, statusCode, data...)
}

// Loader method sets the template loader for this context. This should be done before calling
// `ctx.Render`.
func (ctx *Context) Loader(name string) *Context {
	ctx.templateLoader = name
	return ctx
}

func (ctx *Context) getRawXSRFToken() string {
	token, err := ctx.Cookie("_xsrf")
	if err != nil {
		return ""
	}
	return token
}

func (ctx *Context) checkXSRFToken() bool {
	token := ctx.Request.FormValue("xsrf_token")
	if token == "" {
		return false
	}
	_, tokenA, _ := decodeXSRFToken(token)
	_, tokenB, _ := decodeXSRFToken(ctx.getRawXSRFToken())
	return compareToken(tokenA, tokenB)
}

func (ctx *Context) xsrfToken() string {
	maskedToken := ctx.getRawXSRFToken()
	if maskedToken == "" {
		maskedToken = newXSRFToken()
		ctx.SetCookie("_xsrf", maskedToken, 3600)
	}
	_, tokenBytes, err := decodeXSRFToken(maskedToken)
	if err != nil {
		maskedToken = newXSRFToken()
		ctx.SetCookie("_xsrf", maskedToken, 3600)
		_, tokenBytes, _ = decodeXSRFToken(maskedToken)
	}
	maskBytes := randomBytes(4)
	maskedTokenBytes := append(maskBytes, websocketMask(maskBytes, tokenBytes)...)
	return hex.EncodeToString(maskedTokenBytes)
}

// Render a template file using the built-in Go template engine.
func (ctx *Context) Render(file string, data ...map[string]interface{}) {
	if ctx.templateLoader == "" {
		panic(fmt.Errorf("Template loader has not been set."))
	}
	var renderData map[string]interface{}
	if len(data) == 0 {
		renderData = make(map[string]interface{})
	} else {
		renderData = data[0]
	}
	renderData["xsrf_token"] = ctx.xsrfToken()
	content, err := ctx.App.View.Render(ctx.templateLoader, file, renderData)
	if err != nil {
		panic(err)
	}
	ctx.Send(content)
}

// RenderFromString renders a input string.
func (ctx *Context) RenderFromString(tplSrc string, data ...map[string]interface{}) {
	var renderData map[string]interface{}
	if len(data) == 0 {
		renderData = make(map[string]interface{})
	} else {
		renderData = data[0]
	}
	renderData["xsrf_token"] = ctx.xsrfToken()
	content, e := ctx.App.View.RenderFromString(ctx.templateLoader, tplSrc, renderData)
	if e != nil {
		panic(e)
	}
	ctx.Send(content)
}
