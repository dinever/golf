.. title:: Golf Web Server
.. highlight:: go

Golf
====

`Golf <http://github.com/dinever/golf>`_ is a fast, simple and lightweight micro-web framework for Go. It comes with powerful features and has no dependencies other than the Go Standard Library.

Installation
------------

::

    go get github.com/dinever/golf

Hello World
-----------

Here is a simple "Hello World!" application using Golf::

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


Documentation
-------------

.. toctree::
   :titlesonly:

   quickstart

LICENSE
-------

`Apache License <http://www.apache.org/licenses/LICENSE-2.0.html>`_
