package Golf

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestWrite(t *testing.T) {
	ctx := makeNewContext("GET", "/foo")
	ctx.Write("hello world")
	if !reflect.DeepEqual(ctx.Body, []byte("hello world")) {
		t.Errorf("Context.Write failed.")
	}
}

func TestAbort(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.Abort(500)
	if w.Code != 500 || !ctx.IsSent {
		t.Errorf("Can not abort a context.")
	}
}

func TestRenderFromString(t *testing.T) {
	cases := []struct {
		src    string
		args   map[string]interface{}
		output string
	}{
		{
			"foo {{.Title}} bar",
			map[string]interface{}{"Title": "Hello World"},
			"foo Hello World bar",
		},
	}

	for _, c := range cases {
		r := makeTestHTTPRequest(nil, "GET", "/")
		w := httptest.NewRecorder()
		app := New()
		ctx := NewContext(r, w, app)
		ctx.RenderFromString(c.src, c.args)
		ctx.Send()
		if w.Body.String() != c.output {
			t.Errorf("Can not render from string correctly: %v != %v", w.Body.String(), c.output)
		}
	}
}

func TestJSON(t *testing.T) {
	cases := []struct {
		input  map[string]interface{}
		output string
	}{
		{
			map[string]interface{}{"status": "success", "code": 200},
			`{"code":200,"status":"success"}`,
		},
	}

	for _, c := range cases {
		r := makeTestHTTPRequest(nil, "GET", "/")
		w := httptest.NewRecorder()
		app := New()
		ctx := NewContext(r, w, app)
		ctx.JSON(c.input)
		ctx.Send()
		if w.Body.String() != c.output {
			t.Errorf("Can not return JSON correctly: %v != %v", w.Body.String(), c.output)
		}
		if w.HeaderMap.Get("Content-Type") != `application/json` {
			t.Errorf("Content-Type didn't set properly when calling Context.JSON.")
		}
	}
}
