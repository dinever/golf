package Golf

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func makeTestHttpRequest(body io.Reader) *http.Request {
	req, err := http.NewRequest("GET", "/foo/bar", body)
	if err != nil {
		return nil
	}
	return req
}

func TestContextCreate(t *testing.T) {

	r := makeTestHttpRequest(nil)
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)

	if ctx == nil {
		t.Errorf("Can not create context.")
	}
}
