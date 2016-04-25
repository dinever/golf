package main

import (
	"github.com/dinever/golf"
	"os"
)

func versionHandler(ctx *golf.Context) {
	appVersion, _ := ctx.App.Config.GetString("APP/VERSION", "unavailable")
	ctx.Send(appVersion)
}

func homeHandler(ctx *golf.Context) {
	apiKey, _ := ctx.App.Config.GetString("API_KEY", "0")
	ctx.Send(apiKey)
}

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	app := golf.New()
	app.Get("/version", versionHandler)
	app.Get("/", homeHandler)
	app.Config, err = golf.ConfigFromJSON(file)
	app.Run(":9000")
}
