<a href="http://golf.readme.io"><img width=500px src="/golf-logo-blue.png"></img></a>

[![GoDoc](http://img.shields.io/badge/golf-documentation-blue.svg?style=flat-square)](http://golf.readme.io/docs)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/dinever/golf/master/LICENSE) 
[![Build Status](http://img.shields.io/travis/dinever/golf.svg?style=flat-square)](https://travis-ci.org/dinever/golf) 
[![Coverage Status](http://img.shields.io/coveralls/dinever/golf.svg?style=flat-square)](https://coveralls.io/r/dinever/golf?branch=master)

A web framework in Go.

Homepage: [golf.readme.io](https://golf.readme.io/)

## Installation

    go get github.com/dinever/golf

## Hello World

```go
package main

import "github.com/dinever/golf"

func mainHandler(ctx *Golf.Context) {
  ctx.Write("Hello World!")
}

func pageHandler(ctx *Golf.Context) {
  page, err := ctx.Param("page")
  if err != nil {
    ctx.Abort(500)
    return
  }
  ctx.Write("Page: " + page)
}

func main() {
  app := Golf.New()
  app.Get("/", mainHandler)
  app.Get("/p/:page/", pageHandler)
  app.Run(":9000")
}
```

The website will be available at http://localhost:9000.

## Docs

[golf.readme.io/docs](https://golf.readme.io/docs)

## License

[MIT License](/LICENSE)
