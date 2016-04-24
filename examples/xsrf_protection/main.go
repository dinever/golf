package main

import (
	"github.com/dinever/golf"
)

func mainHandler(ctx *golf.Context) {
	ctx.Send("Hello World!")
}

func loginHandler(ctx *golf.Context) {
	ctx.Loader("default").Render("login.html", make(map[string]interface{}))
}

func loginHandlerPost(ctx *golf.Context) {
	ctx.Send("Hi, " + ctx.Request.FormValue("name"))
}

func main() {
	app := golf.New()
	app.Use(golf.XSRFProtectionMiddleware)
	app.View.SetTemplateLoader("default", ".")

	app.Get("/", mainHandler)
	app.Post("/login", loginHandlerPost)
	app.Get("/login", loginHandler)

	app.Run(":9000")
}
