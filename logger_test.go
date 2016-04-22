package golf

import (
	"bytes"
	"testing"
	"regexp"
)

func assertContains(t *testing.T, content string, query string) {
	re := regexp.MustCompile(query)
	if re.FindString(content) == "" {
		t.Errorf("Not contain: %v in %v", query, content)
	}
}

func TestLogger(t *testing.T) {
	buffer := new(bytes.Buffer)
	app := New()
	app.Use(LoggingMiddleware(buffer))
	app.Get("/example", func(c *Context) {})
	app.Post("/example", func(c *Context) {})
	app.Put("/example", func(c *Context) {})
	app.Delete("/example", func(c *Context) {})
	app.Patch("/example", func(c *Context) {})
	app.Head("/example", func(c *Context) {})
	app.Options("/example", func(c *Context) {})

	_, _, r, w := makeTestContext("GET", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "GET")
	assertContains(t, buffer.String(), "/example")

	// I wrote these first (extending the above) but then realized they are more
	// like integration tests because they test the whole logging process rather
	// than individual functions.  Im not sure where these should go.

	_, _, r, w = makeTestContext("POST", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "POST")
	assertContains(t, buffer.String(), "/example")

	_, _, r, w = makeTestContext("PUT", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "PUT")
	assertContains(t, buffer.String(), "/example")

	_, _, r, w = makeTestContext("DELETE", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "DELETE")
	assertContains(t, buffer.String(), "/example")

	_, _, r, w = makeTestContext("PATCH", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "PATCH")
	assertContains(t, buffer.String(), "/example")

	_, _, r, w = makeTestContext("HEAD", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "HEAD")
	assertContains(t, buffer.String(), "/example")

	_, _, r, w = makeTestContext("OPTIONS", "/example")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "200")
	assertContains(t, buffer.String(), "OPTIONS")
	assertContains(t, buffer.String(), "/example")

	_, _, r, w = makeTestContext("GET", "/notfound")
	app.ServeHTTP(w, r)
	assertContains(t, buffer.String(), "404")
	assertContains(t, buffer.String(), "GET")
	assertContains(t, buffer.String(), "/notfound")
}
