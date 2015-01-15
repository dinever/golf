package Golf

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type Application struct {
	router       *Router
	staticRouter map[string][]string
	view         *View
	Config       *Config
}

func New() *Application {
	app := new(Application)
	app.router = NewRouter()
	app.staticRouter = make(map[string][]string)
	app.view = NewView("")
	app.Config = NewConfig(app)
	return app
}

func (app *Application) handler(req *Request, res *Response) {
	for prefix, staticPathSlice := range app.staticRouter {
		if strings.HasPrefix(req.URL.Path, prefix) {
			for _, staticPath := range staticPathSlice {
				filePath := path.Join(staticPath, req.URL.Path[len(prefix):])
				_, err := os.Stat(filePath)
				if err == nil {
					staticHandler(req, res, filePath)
					return
				}
			}
		}
		notFoundHandler(req, res)
	}

	var (
		params  map[string]string
		handler Handler
	)
	params, handler = app.router.match(req.URL.Path, req.Method)
	if params != nil && handler != nil {
		res.StatusCode = 200
		req.Params = params
		handler(req, res)
	} else {
		notFoundHandler(req, res)
	}
	res.Write(res.Body)
}

func notFoundHandler(req *Request, res *Response) {
	res.StatusCode = 404
	res.Send("404")
}

func staticHandler(req *Request, res *Response, filePath string) {
	http.ServeFile(res.ResponseWriter, req.Request, filePath)
}

func (app *Application) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	request := NewRequest(req)
	response := NewResponse(res, app)
	chain.Final(app.handler)(request, response)
}

func (app *Application) Run(port string) {
	e := http.ListenAndServe(port, app)
	panic(e)
}

func (app *Application) Static(url string, path string) {
	url = strings.TrimRight(url, "/")
	app.staticRouter[url] = append(app.staticRouter[url], path)
}

func (app *Application) Get(pattern string, handler Handler) {
	app.router.Get(pattern, handler)
}

func (app *Application) Post(pattern string, handler Handler) {
	app.router.Post(pattern, handler)
}

func (app *Application) Put(pattern string, handler Handler) {
	app.router.Put(pattern, handler)
}

func (app *Application) Delete(pattern string, handler Handler) {
	app.router.Delete(pattern, handler)
}
