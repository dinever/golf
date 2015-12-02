.. title:: Golf Web Server
.. highlight:: go

Golf
====

`Golf <http://github.com/dinever/golf>`_ is a Go web framework.

Installation
------------

::

    go get github.com/dinever/golf

Hello World
-----------

Here is a simple "Hello World!" application using Golf::

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


Documentation
-------------

.. toctree::
   :titlesonly:

   request
   response

LICENSE
-------

`Apache License <http://www.apache.org/licenses/LICENSE-2.0.html>`_
