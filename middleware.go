package Golf

import (
	"errors"
	"log"
	"time"
)

type middlewareHandler func(next Handler) Handler

var defaultMiddlewares = []middlewareHandler{LoggingMiddleware, RecoverMiddleware, XSRFProtectionMiddleware}

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
func (c Chain) Final(fn Handler) Handler {
	final := fn
	for i := len(c.middlewareHandlers) - 1; i >= 0; i-- {
		final = c.middlewareHandlers[i](final)
	}
	return final
}

// Append a middleware to the middleware chain
func (c *Chain) Append(fn middlewareHandler) {
	c.middlewareHandlers = append(c.middlewareHandlers, fn)
}

// LoggingMiddleware is the built-in middleware for logging.
func LoggingMiddleware(next Handler) Handler {
	fn := func(ctx *Context) {
		t1 := time.Now()
		next(ctx)
		t2 := time.Now()
		log.Printf("[%s] %q %v %v\n", ctx.Request.Method, ctx.Request.URL.String(), ctx.StatusCode, t2.Sub(t1))
	}
	return fn
}

// RecoverMiddleware is the built-in middleware for recovering from errors.
func RecoverMiddleware(next Handler) Handler {
	fn := func(ctx *Context) {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("Unknown panic")
				}
				log.Printf("Panic: %+v", err.Error())
				ctx.App.handleError(ctx, 500)
			}
		}()
		next(ctx)
	}
	return fn
}

// XSRFProtectionMiddleware is the built-in middleware for XSRF protection.
func XSRFProtectionMiddleware(next Handler) Handler {
	fn := func(ctx *Context) {
		xsrfEnabled, _ := ctx.App.Config.GetBool("xsrf_cookies", false)
		if xsrfEnabled && (ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "DELETE") {
			if !checkXSRFToken(ctx) {
				ctx.App.handleError(ctx, 403)
				return
			}
		}
		next(ctx)
	}
	return fn
}
