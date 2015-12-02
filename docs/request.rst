Request
============

``Request`` is a wrapper of Go's built-in ``http.Request``.

func (Request) Query
************************

``func (req *Request) Query(key string, index ...int) (string, error)``

``Query`` returns query data by the query key.::

    // GET /search?q=foo
    val, err := req.Query("q")
    // var == "foo", err == nil

You can indicate an index for multiple query with the same key.::

    // GET /search?q=foo&q=bar
    val, err := req.Query("q", 0)
    // var == "foo", err == nil
    val, err := req.Query("q", 1)
    // var == "bar", err == nil

If there is no query matching the key, ``Query`` returns "".::

    // GET /search?q=foo
    val, err := req.Query("m")
    // var == "", err != nil

func (Request) Cookie
*********************

``func (req *Request) Cookie(key string) string``

``Cookie`` returns the cookie value from the request based on the given key::

    // Cookie: name=foo
    req.Cookie("name")
    // "foo"


func (Request) Protocol
***********************

``func (req *Request) Protocol() string``

``Protocl`` returns the request protocol string.::

    req.Protocol()
    // "HTTP/1.1"
