<a href="http://golf.readme.io"><img width=50% src="/golf-logo.png"></img></a>

[![GoDoc](http://img.shields.io/badge/golf-documentation-blue.svg?style=flat-square)](http://golf.readme.io/docs)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/dinever/golf/master/LICENSE) 
[![Build Status](http://img.shields.io/travis/dinever/golf.svg?style=flat-square)](https://travis-ci.org/dinever/golf) 
[![Build Status](https://goreportcard.com/badge/github.com/dinever/golf?style=flat-square)](https://travis-ci.org/dinever/golf) 
[![Coverage Status](http://img.shields.io/coveralls/dinever/golf.svg?style=flat-square)](https://coveralls.io/r/dinever/golf?branch=master)

A fast, simple and lightweight micro-web framework for Go, comes with powerful features and has no dependencies other than the Go Standard Library.

Homepage: [golf.readme.io](https://golf.readme.io/)

## Installation

    go get github.com/dinever/golf

## Features

1. No allocation during routing and parameter retrieve.
1. Dead simple template inheritance with `extends` and `include` helper comes out of box.

    **layout.html**
    
    ```html
    <h1>Hello World</h1>
    {{ template "body" }}
    {{ include "sidebar.html" }}
    ```
    
    **index.html**

    ```jinja2
    {{ extends "layout.html" }}
    
    {{ define "body"}}
    <p>Main content</p>
    {{ end }}
    ```
    
    **sidebar.html**
    
    ```jinja2
    <p>Sidebar content</p>
    ```
1. Built-in XSRF and Session support.
1. Powerful middleware chain.
1. Configuration from JSON file.

## Hello World

```go
package main

import "github.com/dinever/golf"

func mainHandler(ctx *golf.Context) {
  ctx.Send("Hello World!")
}

func pageHandler(ctx *golf.Context) {
  ctx.Send("Page: " + ctx.Param("page"))
}

func main() {
  app := golf.New()
  app.Get("/", mainHandler)
  app.Get("/p/:page/", pageHandler)
  app.Run(":9000")
}
```

The website will be available at http://localhost:9000.

## Benchmark

The following chart shows the benchmark performance of Golf compared with others.

![Golf benchmark](https://cloud.githubusercontent.com/assets/1311594/14748305/fcbdc216-0886-11e6-90a4-231e78acfb60.png)

For more information, please see [BENCHMARKING.md](BENCHMARKING.md)

## Docs

[golf.readme.io/docs](https://golf.readme.io/docs)

## License

[MIT License](/LICENSE)
