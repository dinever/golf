# Golf

A web framework in Go.

Homepage: [golf.readthedocs.org](http://crotal.org)

## Installation

    go get github.com/dinever/golf

## Hello World

    package main

    import "github.com/dinever/golf"

    func helloWorldHandler(req Golf.Request, res Golf.Response) {
        res.Send("Hello World!")
    }

    func main() {
        app := Golf.New()
        app.Get("/", helloWorldHandler)
        app.Run(":5693")
    }

The website will be available at http://localhost:5693.

##Documents

See [golf.readthedocs.org](http://crotal.org).

##License

[Apache License](http://www.apache.org/licenses/LICENSE-2.0.html)
