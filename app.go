package Golf

import (
	"net/http"
	"os"
	"path"
	"strings"
)

// Application is an abstraction of a Golf application, can be used for
// configuration, etc.
type Application struct {
	router *router

	// A map of string slices as value to indicate the static files.
	staticRouter map[string][]string

	// The View model of the application. View handles the templating and page
	// rendering.
	View *View

	// Config provides configuration management.
	Config *Config

	SessionManager SessionManager

	// NotFoundHandler handles requests when no route is matched.
	NotFoundHandler Handler

	// MiddlewareChain is the default middlewares that Golf uses.
	MiddlewareChain *Chain

	errorHandler map[int]ErrorHandlerType

	// The default error handler, if the corresponding error code is not specified
	// in the `errorHandler` map, this handler will be called.
	DefaultErrorHandler ErrorHandlerType
}

// New is used for creating a new Golf Application instance.
func New() *Application {
	app := new(Application)
	app.router = newRouter()
	app.staticRouter = make(map[string][]string)
	app.View = NewView()
	app.Config = NewConfig(app)
	// debug, _ := app.Config.GetBool("debug", false)
	app.errorHandler = make(map[int]ErrorHandlerType)
	app.MiddlewareChain = NewChain(defaultMiddlewares...)
	app.DefaultErrorHandler = defaultErrorHandler
	return app
}

// First search if any of the static route matches the request.
// If not, look up the URL in the router.
func (app *Application) handler(ctx *Context) {
	for prefix, staticPathSlice := range app.staticRouter {
		if strings.HasPrefix(ctx.Request.URL.Path, prefix) {
			for _, staticPath := range staticPathSlice {
				filePath := path.Join(staticPath, ctx.Request.URL.Path[len(prefix):])
				fileInfo, err := os.Stat(filePath)
				if err == nil && !fileInfo.IsDir() {
					staticHandler(ctx, filePath)
					return
				}
			}
		}
	}

	params, handler := app.router.match(ctx.Request.URL.Path, ctx.Request.Method)
	if handler != nil {
		ctx.Params = params
		handler(ctx)
	} else {
		app.handleError(ctx, 404)
	}
	ctx.Send()
}

// Serve a static file
func staticHandler(ctx *Context, filePath string) {
	http.ServeFile(ctx.Response, ctx.Request, filePath)
}

// Basic entrance of an `http.ResponseWriter` and an `http.Request`.
func (app *Application) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := NewContext(req, res, app)
	app.MiddlewareChain.Final(app.handler)(ctx)
}

// Run the Golf Application.
func (app *Application) Run(addr string) {
	err := http.ListenAndServe(addr, app)
	if err != nil {
		panic(err)
	}
}

// RunTLS runs the app with TLS support.
func (app *Application) RunTLS(addr, certFile, keyFile string) {
	err := http.ListenAndServeTLS(addr, certFile, keyFile, app)
	if err != nil {
		panic(err)
	}
}

// Static is used for registering a static folder
func (app *Application) Static(url string, path string) {
	url = strings.TrimRight(url, "/")
	app.staticRouter[url] = append(app.staticRouter[url], path)
}

// Get method is used for registering a Get method route
func (app *Application) Get(pattern string, handler Handler) {
	app.router.get(pattern, handler)
}

// Post method is used for registering a Post method route
func (app *Application) Post(pattern string, handler Handler) {
	app.router.post(pattern, handler)
}

// Put method is used for registering a Put method route
func (app *Application) Put(pattern string, handler Handler) {
	app.router.put(pattern, handler)
}

// Delete method is used for registering a Delete method route
func (app *Application) Delete(pattern string, handler Handler) {
	app.router.delete(pattern, handler)
}

// Error method is used for registering an handler for a specified HTTP error code.
func (app *Application) Error(statusCode int, handler ErrorHandlerType) {
	app.errorHandler[statusCode] = handler
}

// Handles a HTTP Error, if there is a corresponding handler set in the map
// `errorHandler`, then call it. Otherwise call the `defaultErrorHandler`.
func (app *Application) handleError(ctx *Context, statusCode int, data ...map[string]interface{}) {
	ctx.StatusCode = statusCode
	handler, ok := app.errorHandler[ctx.StatusCode]
	if !ok {
		defaultErrorHandler(ctx, data...)
		return
	}
	handler(ctx)
}
