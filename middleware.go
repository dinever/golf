package Golf

import (
	"errors"
	"log"
	"time"
)

type middlewareHandler func(next Handler) Handler

var defaultMiddlewares = []middlewareHandler{LoggingMiddleware, RecoverMiddleware}

// A chain of middlewares.
type Chain struct {
	middlewareHandlers []middlewareHandler
}

func NewChain(handlerArray ...middlewareHandler) *Chain {
	c := new(Chain)
	c.middlewareHandlers = handlerArray
	return c
}

// Indicating a final Handler, chain the multiple middlewares together with the
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

func LoggingMiddleware(next Handler) Handler {
	fn := func(ctx *Context) {
		t1 := time.Now()
		next(ctx)
		t2 := time.Now()
		log.Printf("[%s] %q %v %v\n", ctx.Request.Method, ctx.Request.URL.String(), ctx.StatusCode, t2.Sub(t1))
	}
	return fn
}

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
