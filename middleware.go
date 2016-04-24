package golf

import (
	"fmt"
	"io"
	"log"
	"net/http/httputil"
	"time"
)

// MiddlewareHandlerFunc defines the middleware function type that Golf uses.
type MiddlewareHandlerFunc func(next HandlerFunc) HandlerFunc

// Chain contains a sequence of middlewares.
type Chain struct {
	middlewareHandlers []MiddlewareHandlerFunc
}

// NewChain Creates a new middleware chain.
func NewChain(handlerArray ...MiddlewareHandlerFunc) *Chain {
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
func (c *Chain) Append(fn MiddlewareHandlerFunc) {
	c.middlewareHandlers = append(c.middlewareHandlers, fn)
}

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

// LoggingMiddleware is the built-in middleware for logging.
// This method is referred from https://github.com/gin-gonic/gin/blob/develop/logger.go#L46
func LoggingMiddleware(out io.Writer) MiddlewareHandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		fn := func(ctx *Context) {
			start := time.Now()
			path := ctx.Request.URL.Path

			next(ctx)

			end := time.Now()
			latency := end.Sub(start)

			clientIP := ctx.ClientIP()
			method := ctx.Request.Method
			statusCode := ctx.statusCode
			statusColor := colorForStatus(statusCode)
			methodColor := colorForMethod(method)

			fmt.Fprintf(out, "%v |%s %3d %s| %13v | %s |%s  %s %-7s %s\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, reset,
				latency,
				clientIP,
				methodColor, reset, method,
				path,
			)
		}
		return fn
	}
}

func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return white
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return blue
	case "POST":
		return cyan
	case "PUT":
		return yellow
	case "DELETE":
		return red
	case "PATCH":
		return green
	case "HEAD":
		return magenta
	case "OPTIONS":
		return white
	default:
		return reset
	}
}

// XSRFProtectionMiddleware is the built-in middleware for XSRF protection.
func XSRFProtectionMiddleware(next HandlerFunc) HandlerFunc {
	fn := func(ctx *Context) {
		if ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "DELETE" {
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
				ctx.statusCode = 500
				ctx.Abort(500, map[string]interface{}{
					"Code":        ctx.statusCode,
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
