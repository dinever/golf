package Yafg

import (
  "testing"
  "reflect"
)

func handler() {}

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
    params []string
  }{
    {"/:year/:mongth/:day/",
      `^/([\w-%]+)/([\w-%]+)/([\w-%]+)/$`,
      []string{"year", "month", "day"}},
  }

  for _, c := range cases {
    route := newRoute(RouterMethodGet, c.in, handler)
    if route.regex.String() != c.regex {
      t.Errorf("regex == %q, want %q", c.in, route.regex.String(), c.regex)
    }
    if reflect.DeepEqual(route.params, c.params) {
      t.Errorf("regex == %q, want %q", c.in, route.regex.String(), c.regex)
    }
  }
}
