package golf

import (
	"log"
	"net/http/httputil"
	"time"
)

type middlewareHandler func(next HandlerFunc) HandlerFunc

var defaultMiddlewares = []middlewareHandler{LoggingMiddleware, RecoverMiddleware, XSRFProtectionMiddleware, SessionMiddleware}

// Chain contains a sequence of middlewares.
type Chain struct {
	middlewareHandlers []middlewareHandler
}

// NewChain Creates a new middleware chain.
func NewChain(handlerArray ...middlewareHandler) *Chain {
	c := new(Chain)
	c.middlewareHandlers = handlerArray
	return c
}

// Final indicates a final Handler, chain the multiple middlewares together with the
// handler, and return them together as a handler.
func (c Chain) Final(fn HandlerFunc) HandlerFunc {
	for i := len(c.middlewareHandlers) - 1; i >= 0; i-- {
		fn = c.middlewareHandlers[i](fn)
	}
	return fn
}

// Append a middleware to the middleware chain
func (c *Chain) Append(fn middlewareHandler) {
	c.middlewareHandlers = append(c.middlewareHandlers, fn)
}

// LoggingMiddleware is the built-in middleware for logging.
func LoggingMiddleware(next HandlerFunc) HandlerFunc {
	fn := func(ctx *Context) {
		t1 := time.Now()
		next(ctx)
		t2 := time.Now()
		log.Printf("[%s] %q %v %v\n", ctx.Request.Method, ctx.Request.URL.String(), ctx.StatusCode, t2.Sub(t1))
	}
	return fn
}

// XSRFProtectionMiddleware is the built-in middleware for XSRF protection.
func XSRFProtectionMiddleware(next HandlerFunc) HandlerFunc {
	fn := func(ctx *Context) {
		xsrfEnabled, _ := ctx.App.Config.GetBool("xsrf_cookies", false)
		if xsrfEnabled && (ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "DELETE") {
			if !ctx.checkXSRFToken() {
				ctx.Abort(403)
				return
			}
		}
		next(ctx)
	}
	return fn
}

// RecoverMiddleware is the built-in middleware for recovering from errors.
func RecoverMiddleware(next HandlerFunc) HandlerFunc {
	fn := func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				e := NewError(err)
				httpRequest, _ := httputil.DumpRequest(ctx.Request, true)
				log.Printf("[Recovery] panic recovered:\n%s\n%s\n%s", string(httpRequest), err, e.StackTraceString())
				ctx.StatusCode = 500
				ctx.Abort(500, map[string]interface{}{
					"Code":        ctx.StatusCode,
					"Title":       "Internal Server Error",
					"HTTPRequest": string(httpRequest),
					"Message":     e.Error(),
					"StackTrace":  e.Stack,
				})
			}
		}()
		next(ctx)
	}
	return fn
}

// SessionMiddleware handles session of the request
func SessionMiddleware(next HandlerFunc) HandlerFunc {
	fn := func(ctx *Context) {
		if ctx.App.SessionManager != nil {
			ctx.retrieveSession()
		}
		next(ctx)
	}
	return fn
}
