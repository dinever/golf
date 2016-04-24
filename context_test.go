package golf

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("Not equal: %v (expected) != %v (actual)", expected, actual)
	}
}

func assertNotEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected == actual {
		t.Errorf("Equal: %v (expected) == %v (actual)", expected, actual)
	}
}

func assertDeepEqual(t *testing.T, expected interface{}, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Not equal: %v (expected) != %v (actual)", expected, actual)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assertError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Should have raised an error")
	}
}

func makeTestHTTPRequest(body io.Reader, method, url string) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil
	}
	return req
}

func makeTestContext(method, url string) (*Context, *Application, *http.Request, *httptest.ResponseRecorder) {
	r := makeTestHTTPRequest(nil, method, url)
	w := httptest.NewRecorder()
	app := New()
	return NewContext(r, w, app), app, r, w
}

func TestContextCreate(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/foo/bar/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	assertNotEqual(t, ctx, nil)
}

func TestParam(t *testing.T) {
	_, app, r, w := makeTestContext("POST", "/foo/")
	app.middlewareChain = NewChain()
	app.Post("/:page/", func(ctx *Context) {
		assertEqual(t, ctx.Param("page"), "foo")
		ctx.Send("success")
	})
	app.ServeHTTP(w, r)
}

func TestParamWithMultipleParameters(t *testing.T) {
	_, app, r, w := makeTestContext("POST", "/dinever/golf/")
	app.middlewareChain = NewChain()
	app.Post("/:user/:repo/", func(ctx *Context) {
		assertEqual(t, ctx.Param("user"), "dinever")
		assertEqual(t, ctx.Param("repo"), "golf")
		assertEqual(t, ctx.Param("org"), "")
		ctx.Send("success")
	})
	app.ServeHTTP(w, r)
}

func TestCookieSet(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/foo/bar/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.SetCookie("foo", "bar", 0)
	assertEqual(t, w.HeaderMap.Get("Set-Cookie"), `foo=bar; Path=/`)
}

func TestCookieSetWithExpire(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/foo/bar/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.SetCookie("foo", "bar", 3600)
	rawCookie := w.HeaderMap.Get("Set-Cookie")
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookie)
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	if err == nil {
		cookies := req.Cookies()
		cookie := cookies[3]
		assertEqual(t, cookie.Value, "3600")
	}
}

func TestSetHeader(t *testing.T) {
	ctx, _, _, w := makeTestContext("GET", "/foo/bar/")
	ctx.SetHeader("foo", "bar")
	assertEqual(t, w.Header().Get("foo"), "bar")
}

func TestAddHeader(t *testing.T) {
	ctx, _, _, w := makeTestContext("GET", "/foo/bar/")
	ctx.AddHeader("foo", "foo")
	ctx.AddHeader("foo", "bar")
	assertDeepEqual(t, w.HeaderMap["Foo"], []string{"foo", "bar"})
}

func TestGetHeader(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/foo/bar/")
	r.Header.Set("foo", "bar")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	assertEqual(t, ctx.Header("foo"), "bar")
}

func TestSendStatus(t *testing.T) {
	ctx, _, _, w := makeTestContext("GET", "/")
	ctx.SendStatus(500)
	assertEqual(t, w.Code, 500)
}

func TestSessionWithInvalidSid(t *testing.T) {
	ctx, app, r, w := makeTestContext("GET", "/foo/bar/")
	app.SessionManager = NewMemorySessionManager()
	ctx.retrieveSession()
	assertNotEqual(t, ctx.Session, nil)
	r.AddCookie(&http.Cookie{Name: "sid", Value: "abc"})
	ctx = NewContext(r, w, app)
	ctx.retrieveSession()
	assertNotEqual(t, ctx.Session, nil)
}

func TestSession(t *testing.T) {
	_, app, r, w := makeTestContext("GET", "/foo/")
	app.SessionManager = NewMemorySessionManager()
	app.Use(SessionMiddleware)
	var (
		firstSid string
	)
	app.Get("/foo/", func(ctx *Context) {
		if ctx.Session == nil {
			t.Errorf("Could not retrieve session.")
		}
		firstSid = ctx.Session.SessionID()
		ctx.Send("success")
	})
	app.ServeHTTP(w, r)
	app.Get("/bar/", func(ctx *Context) {
		if ctx.Session.SessionID() != firstSid {
			t.Errorf("Could not retrieve correct session from the same user.")
		}
		ctx.Send("success")
	})
	_, _, r, w = makeTestContext("GET", "/bar/")
	r.AddCookie(&http.Cookie{Name: "sid", Value: firstSid})
	app.ServeHTTP(w, r)
}

func TestXSRFProtectionWithoutCookie(t *testing.T) {
	ctx, app, _, _ := makeTestContext("GET", "/foo/bar/")
	app.Config.Set("xsrf_cookies", true)
	assertEqual(t, ctx.getRawXSRFToken(), "")
}

func TestXSRFProtectionDisabled(t *testing.T) {
	_, app, r, w := makeTestContext("POST", "/foo/bar/")
	app.middlewareChain = NewChain(XSRFProtectionMiddleware)
	app.Post("/foo/bar/", func(ctx *Context) {
		ctx.Send("success")
	})
	app.ServeHTTP(w, r)

	assertNotEqual(t, w.Code, 403)
	assertEqual(t, w.Body.String(), "success")
}

func TestXSRFProtection(t *testing.T) {
	_, app, r, w := makeTestContext("GET", "/login/")
	app.Config.Set("xsrf_cookies", true)
	app.middlewareChain = NewChain(XSRFProtectionMiddleware)
	var expectedToken string
	app.Get("/login/", func(ctx *Context) {
		expectedToken = ctx.xsrfToken()
		ctx.Send("success")
	})
	app.Post("/login/", func(ctx *Context) {
		ctx.Send("success")
	})
	app.ServeHTTP(w, r)

	_, tokenBytes, _ := decodeXSRFToken(expectedToken)
	maskBytes := randomBytes(4)
	maskedTokenBytes := append(maskBytes, websocketMask(maskBytes, tokenBytes)...)
	maskedToken := hex.EncodeToString(maskedTokenBytes)

	_, _, r, w = makeTestContext("POST", "/login/")
	r.AddCookie(&http.Cookie{Name: "_xsrf", Value: expectedToken})
	r.Form.Add("xsrf_token", maskedToken)
	app.ServeHTTP(w, r)
	assertEqual(t, w.Code, 200)
}

func TestTemplateLoader(t *testing.T) {
	ctx, _, _, _ := makeTestContext("GET", "/")
	ctx.Loader("admin")
	assertEqual(t, ctx.templateLoader, "admin")
}

func TestQuery(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/search?q=foo&p=bar")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	q, err := ctx.Query("q")
	if err != nil {
		t.Errorf("Could not retrieve a query.")
	} else {
		if q != "foo" {
			t.Errorf("Could not retrieve the correct query `q`.")
		}
	}
	p, err := ctx.Query("p")
	assertNoError(t, err)
	assertEqual(t, p, "bar")
}

func TestQueries(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/search?myarray=value1&myarray=value2&myarray=value3")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	q, err := ctx.Query("myarray", 2)
	assertNoError(t, err)
	assertEqual(t, q, "value3")
}

func TestQueryNotFound(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/search?myarray=value1&myarray=value2&myarray=value3")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	q, err := ctx.Query("query")
	assertError(t, err)
	assertEqual(t, q, "")
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
	assertEqual(t, w.Header().Get("Location"), `/foo`)
	assertEqual(t, w.Code, 301)
}

func TestSendString(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.Send("hello world")
	assertDeepEqual(t, w.Body.Bytes(), []byte("hello world"))
}

func TestSendByteSlice(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.Send([]byte("hello world"))
	assertDeepEqual(t, w.Body.Bytes(), []byte("hello world"))
}

func TestSendBuffer(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	buf := bytes.NewBuffer([]byte("hello world"))
	ctx.Send(buf)
	assertDeepEqual(t, w.Body.Bytes(), []byte("hello world"))
}

func TestAbort(t *testing.T) {
	r := makeTestHTTPRequest(nil, "GET", "/")
	w := httptest.NewRecorder()
	app := New()
	ctx := NewContext(r, w, app)
	ctx.Abort(500)
	if w.Code != 500 || !ctx.IsSent {
		t.Errorf("Could not abort a context.")
	}
}

func TestContextClientIP(t *testing.T) {
	ctx := makeNewContext("POST", "/")
	ctx.Request.Header.Set("X-Real-IP", " 10.10.10.10  ")
	ctx.Request.Header.Set("X-Forwarded-For", "  20.20.20.20, 30.30.30.30")
	ctx.Request.RemoteAddr = "  40.40.40.40:42123 "

	assertEqual(t, ctx.ClientIP(), "10.10.10.10")
	assertEqual(t, ctx.ClientIP(), "10.10.10.10")

	ctx.Request.Header.Del("X-Real-IP")
	assertEqual(t, ctx.ClientIP(), "20.20.20.20")

	ctx.Request.Header.Set("X-Forwarded-For", "30.30.30.30  ")
	assertEqual(t, ctx.ClientIP(), "30.30.30.30")

	ctx.Request.Header.Del("X-Forwarded-For")
	assertEqual(t, ctx.ClientIP(), "40.40.40.40")
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
		assertEqual(t, w.Body.String(), c.output)
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
		assertEqual(t, w.Body.String(), c.output)
		assertEqual(t, w.HeaderMap.Get("Content-Type"), `application/json`)
	}
}
