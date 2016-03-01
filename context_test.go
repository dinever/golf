package Golf

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTestHTTPRequest(body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", "/foo/bar", body)
	if err != nil {
		return nil
	}
	return req
}

func TestContextCreate(t *testing.T) {
	r := makeTestHTTPRequest(nil)
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	if ctx == nil {
		t.Errorf("Can not create context.")
	}
}

func TestCookieSet(t *testing.T) {
	r := makeTestHTTPRequest(nil)
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.SetCookie("foo", "bar", 0)
	ctx.Send()
	if w.HeaderMap.Get("Set-Cookie") != `foo=bar; Path=/` {
		t.Errorf("Cookie test failed: %q != %q", w.HeaderMap.Get("Set-Cookie"), `foo=bar; Path=/`)
	}
}
