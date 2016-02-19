package Golf

import (
	"reflect"
	"testing"
)

func handler(ctx *Context) {}

func TestParsePatternWithOneParam(t *testing.T) {
	cases := []struct {
		in, regex, param string
	}{
		{"/:id/", `^/([\w-%]+)/$`, "id"},
	}

	for _, c := range cases {
		route := newRoute(RouterMethodGet, c.in, handler)
		if route.regex.String() != c.regex {
			t.Errorf("regex of %q  == %q, want %q", c.in, route.regex.String(), c.regex)
		}
		if len(route.params) != 1 {
			t.Errorf("%q is supposed to have 1 parameter", c.in)
		}
		if route.params[0] != "id" {
			t.Errorf("params[0] == %q, want %q", c.in, c.param)
		}
	}
}

func TestParsePatternWithThreeParam(t *testing.T) {
	cases := []struct {
		in, regex string
		params    []string
	}{
		{
			"/:year/:month/:day/",
			`^/([\w-%]+)/([\w-%]+)/([\w-%]+)/$`,
			[]string{"year", "month", "day"},
		},
	}

	for _, c := range cases {
		route := newRoute(RouterMethodGet, c.in, handler)
		if route.regex.String() != c.regex {
			t.Errorf("regex == %q, want %q", c.in, route.regex.String(), c.regex)
		}
		if !reflect.DeepEqual(route.params, c.params) {
			t.Errorf("parameters not match: %v != %v", route.params, c.params)
		}
	}
}

func TestRouterMatch(t *testing.T) {
	router := NewRouter()
	cases := []struct {
		pattern string
		url     string
		params  map[string]string
	}{
		{
			"/:year/:month/:day/",
			"/2015/11/15/",
			map[string]string{"year": "2015", "month": "11", "day": "15"},
		},
		{
			"/user/:id/",
			"/user/foobar/",
			map[string]string{"id": "foobar"},
		},
	}
	for _, c := range cases {
		router.Get(c.pattern, handler)
		params, _ := router.match(c.url, RouterMethodGet)
		if !reflect.DeepEqual(params, c.params) {
			t.Errorf("parameters not match: %v != %v", params, c.params)
		}
	}
}
