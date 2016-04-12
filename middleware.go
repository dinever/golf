package Golf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"runtime"
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

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = []byte{27, 91, 48, 109}
)

// RecoverMiddleware is the built-in middleware for recovering from errors.
func RecoverMiddleware(next Handler) Handler {
	fn := func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				httprequest, _ := httputil.DumpRequest(ctx.Request, false)
				log.Printf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, string(reset))
				ctx.Abort(500)
			}
		}()
		next(ctx)
	}
	return fn
}

// Code initially taken from Gin
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
