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
}

func New() *Application {
	app := new(Application)
	app.router = NewRouter()
	app.staticRouter = make(map[string][]string)
	app.view = NewView("")
	return app
}

func (app *Application) handler(res http.ResponseWriter, req *http.Request) {
	request := *NewRequest(req)
	response := *NewResponse(res, app)

	for prefix, staticPathSlice := range app.staticRouter {
		if strings.HasPrefix(request.URL.Path, prefix) {
			for _, staticPath := range staticPathSlice {
				filePath := path.Join(staticPath, request.URL.Path[len(prefix):])
				_, err := os.Stat(filePath)
				if err == nil {
					staticHandler(request, response, filePath)
					return
				}
			}
			response.Send("404")
		}
	}

	var (
		params  map[string]string
		handler Handler
	)
	params, handler = app.router.match(request.URL.Path, request.Method)
	if params != nil && handler != nil {
		request.Params = params
		handler(request, response)
	} else {
		response.Send("404")
	}
}

func staticHandler(req Request, res Response, filePath string) {
	http.ServeFile(res.ResponseWriter, req.Request, filePath)
}

func (app *Application) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.handler(res, req)
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
