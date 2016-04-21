package golf

import (
	"golang.org/x/net/context"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

type handlerFunc func(c context.Context, w http.ResponseWriter, r *http.Request)

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

	errorHandler map[int]ErrorHandlerType

	// The default error handler, if the corresponding error code is not specified
	// in the `errorHandler` map, this handler will be called.
	DefaultErrorHandler ErrorHandlerType

	pool sync.Pool
}

// New is used for creating a new Golf Application instance.
func New() *Application {
	app := new(Application)
	app.router = newRouter()
	app.staticRouter = make(map[string][]string)
	app.View = NewView()
	app.Config = NewConfig(app)
	app.pool.New = func() interface{} {
		return NewContext()
	}
	// debug, _ := app.Config.GetBool("debug", false)
	app.errorHandler = make(map[int]ErrorHandlerType)
	return app
}

// First search if any of the static route matches the request.
// If not, look up the URL in the router.
func (app *Application) handler(c *Context, w http.ResponseWriter, r *http.Request) {
	for prefix, staticPathSlice := range app.staticRouter {
		if strings.HasPrefix(r.URL.Path, prefix) {
			for _, staticPath := range staticPathSlice {
				filePath := path.Join(staticPath, r.URL.Path[len(prefix):])
				fileInfo, err := os.Stat(filePath)
				if err == nil && !fileInfo.IsDir() {
					http.ServeFile(w, r, filePath)
					return
				}
			}
		}
	}

	handler, params, err := app.router.FindRoute(r.Method, r.URL.Path)
	c.Params = params
	if err != nil {
		//		app.handleError(ctx, 404)
	} else {
		//		ctx.Params = params
		handler(c, w, r)
	}
}

// Basic entrance of an `http.ResponseWriter` and an `http.Request`.
func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTPContext(context.Background(), w, r)
}

func (app *Application) ServeHTTPContext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c := app.pool.Get().(*Context)
	c.parent = ctx

	app.handler(c, w, r)

	c.reset()
	app.pool.Put(c)
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
func (app *Application) Get(pattern string, handler handlerFunc) {
	app.router.AddRoute("GET", pattern, handler)
}

// Post method is used for registering a Post method route
func (app *Application) Post(pattern string, handler handlerFunc) {
	app.router.AddRoute("POST", pattern, handler)
}

// Put method is used for registering a Put method route
func (app *Application) Put(pattern string, handler handlerFunc) {
	app.router.AddRoute("PUT", pattern, handler)
}

// Delete method is used for registering a Delete method route
func (app *Application) Delete(pattern string, handler handlerFunc) {
	app.router.AddRoute("DELETE", pattern, handler)
}

// Error method is used for registering an handler for a specified HTTP error code.
func (app *Application) Error(statusCode int, handler ErrorHandlerType) {
	app.errorHandler[statusCode] = handler
}
