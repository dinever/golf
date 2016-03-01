package Golf

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTestHTTPRequest(body io.Reader, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	return req
}

func TestContextCreate(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/foo/bar/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	if ctx == nil {
		t.Errorf("Can not create context.")
	}
}

func TestCookieSet(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/foo/bar/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.SetCookie("foo", "bar", 0)
	ctx.Send()
	if w.HeaderMap.Get("Set-Cookie") != `foo=bar; Path=/` {
		t.Errorf("Cookie test failed: %q != %q", w.HeaderMap.Get("Set-Cookie"), `foo=bar; Path=/`)
	}
}

func TestQuery(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/search?q=foo&p=bar")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	q, err := ctx.Query("q")
	if err != nil {
		t.Errorf("Can not retrieve a query.")
	} else {
		if q != "foo" {
			t.Errorf("Can not retrieve the correct query `q`.")
		}
	}
	p, err := ctx.Query("p")
	if err != nil {
		t.Errorf("Can not retrieve a query.")
	} else {
		if p != "bar" {
			t.Errorf("Can not retrieve the correct query `p`.")
		}
	}
}

func makeNewContext(method, url string) *Context {
	r := makeTestHTTPRequest(nil, method, url)
	w := httptest.NewRecorder()
	app := New()
	return NewContext(r, w, app)
}

func TestRedirection(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.Redirect("/foo")
	ctx.Send()
	if w.HeaderMap.Get("Location") != `/foo` {
		t.Errorf("Can not perform a 301 redirection.")
	}
}
