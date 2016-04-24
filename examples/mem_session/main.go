package main

import (
	"github.com/dinever/golf"
	"strconv"
)

func mainHandler(ctx *golf.Context) {
	name, err := ctx.Session.Get("name")
	if err != nil {
		ctx.Send("Hello World! Please <a href=\"/login\">log in</a>. Current sessions: " + strconv.Itoa(ctx.App.SessionManager.Count()))
	} else {
		ctx.Send("Hello " + name.(string) + ". Current sessions: " + strconv.Itoa(ctx.App.SessionManager.Count()))
	}
}

func loginHandler(ctx *golf.Context) {
	ctx.Loader("default").Render("login.html", make(map[string]interface{}))
}

func loginHandlerPost(ctx *golf.Context) {
	ctx.Session.Set("name", ctx.Request.FormValue("name"))
	ctx.Send("Hi, " + ctx.Request.FormValue("name"))
}

func main() {
	app := golf.New()
	app.View.SetTemplateLoader("default", ".")
	app.SessionManager = golf.NewMemorySessionManager()
	app.Use(golf.SessionMiddleware)

	app.Get("/", mainHandler)
	app.Post("/login", loginHandlerPost)
	app.Get("/login", loginHandler)

	app.Run(":9000")
}
