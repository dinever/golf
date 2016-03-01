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
	router *Router

	// A map of string slices as value to indicate the static files.
	staticRouter map[string][]string

	// The View model of the application. View handles the templating and page
	// rendering.
	View *View

	// Config provides configuration management.
	Config *Config

	// NotFoundHandler handles requests when no route is matched.
	NotFoundHandler Handler

	// MiddlewareChain is the default middlewares that Golf uses.
	MiddlewareChain *Chain

	errorHandler map[int]Handler

	// The default error handler, if the corresponding error code is not specified
	// in the `errorHandler` map, this handler will be called.
	DefaultErrorHandler Handler
}

// New is used for creating a new Golf Application instance.
func New() *Application {
	app := new(Application)
	app.router = NewRouter()
	app.staticRouter = make(map[string][]string)
	app.View = NewView()
	app.Config = NewConfig(app)
	// debug, _ := app.Config.GetBool("debug", false)
	app.errorHandler = make(map[int]Handler)
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

// Static is used for registering a static folder
func (app *Application) Static(url string, path string) {
	url = strings.TrimRight(url, "/")
	app.staticRouter[url] = append(app.staticRouter[url], path)
}

// Get method is used for registering a Get method route
func (app *Application) Get(pattern string, handler Handler) {
	app.router.Get(pattern, handler)
}

// Post method is used for registering a Post method route
func (app *Application) Post(pattern string, handler Handler) {
	app.router.Post(pattern, handler)
}

// Put method is used for registering a Put method route
func (app *Application) Put(pattern string, handler Handler) {
	app.router.Put(pattern, handler)
}

// Delete method is used for registering a Delete method route
func (app *Application) Delete(pattern string, handler Handler) {
	app.router.Delete(pattern, handler)
}

// Error method is used for registering an handler for a specified HTTP error code.
func (app *Application) Error(statusCode int, handler Handler) {
	app.errorHandler[statusCode] = handler
}

// Handles a HTTP Error, if there is a corresponding handler set in the map
// `errorHandler`, then call it. Otherwise call the `defaultErrorHandler`.
func (app *Application) handleError(ctx *Context, statusCode int) {
	ctx.StatusCode = statusCode
	handler, ok := app.errorHandler[ctx.StatusCode]
	if !ok {
		defaultErrorHandler(ctx)
		return
	}
	handler(ctx)
}
