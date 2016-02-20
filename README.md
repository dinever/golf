# Golf

[![Build Status](https://drone.io/github.com/dinever/golf/status.png)](https://drone.io/github.com/dinever/golf/latest)

A web framework in Go.

Homepage: [golf.readthedocs.org](http://golf.readthedocs.org)

## Installation

    go get github.com/dinever/golf

## Hello World

```go
package main

import "github.com/dinever/golf"

func helloWorldHandler(ctx *Golf.Context) {
  ctx.Write("Hello World!")
}

func main() {
  app := Golf.New()
  app.Get("/", helloWorldHandler)
  app.Run(":9000")
}
```

The website will be available at http://localhost:9000.

##Documents

See [golf.readthedocs.org](http://golf.readthedocs.org).

##License

[Apache License](http://www.apache.org/licenses/LICENSE-2.0.html)
