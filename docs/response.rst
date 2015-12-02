Response
============

.. highlight:: go


``Response`` is a wrapper of Go's built-in ``http.ResponseWriter``.

func (Response) Send
************************

``func (res *Response) Send(str string)``

``Send`` simply send a string as the HTTP response::

    res.Send("Hello World!")

func (Response) Redirect
************************

``func (res *Response) Redirect(url string, code int)``

Redirect to the specified ``url`` with the specified status code ``code``.::

    res.Redirect("/foo", 302)

func (Response) SetCookie
*************************

``func (res *Response) SetCookie(key string, value string, expire int)``

Set cookie ``key`` to ``value``. The cookie will expire after ``expire`` seconds.
