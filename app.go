package Yafg

import (
  "net/http"
)

type Application struct {
  router *Router
}

func New() *Application {
  app := new(Application)
  app.router = NewRouter()
  return app
}

func (app *Application) handler (res http.ResponseWriter, req *http.Request) {
  request := *NewRequest(req)
  response := *NewResponse(res)
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

func (app *Application) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.handler(res, req)
}

func (app *Application) Run(port string) {
  e := http.ListenAndServe(port, app)
	panic(e)
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
