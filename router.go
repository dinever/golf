package Golf

import (
	"regexp"
	"strings"
)

const (
	routerMethodGet    = "GET"
	routerMethodPost   = "POST"
	routerMethodPut    = "PUT"
	routerMethodDelete = "DELETE"
)

type router struct {
	routeSlice []*route
}

// Handler is the type of the handler function that Golf accepts.
type Handler func(ctx *Context)
// ErrorHandler is the type of the error handler function that Golf accepts. 
type ErrorHandler func(ctx *Context, e error)

type route struct {
	method  string
	pattern string
	regex   *regexp.Regexp
	params  []string
	handler Handler
}

func newRouter() *router {
	r := new(router)
	r.routeSlice = make([]*route, 0)
	return r
}

func newRoute(method string, pattern string, handler Handler) *route {
	route := new(route)
	route.pattern = pattern
	route.params = make([]string, 0)
	route.regex, route.params = route.parseURL(pattern)
	route.method = method
	route.handler = handler
	return route
}

func (router *router) get(pattern string, handler Handler) {
	route := newRoute(routerMethodGet, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *router) post(pattern string, handler Handler) {
	route := newRoute(routerMethodPost, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *router) put(pattern string, handler Handler) {
	route := newRoute(routerMethodPut, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *router) delete(pattern string, handler Handler) {
	route := newRoute(routerMethodDelete, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *router) registerRoute(route *route, handler Handler) {
	router.routeSlice = append(router.routeSlice, route)
}

func (router *router) match(url string, method string) (params map[string]string, handler Handler) {
	params = make(map[string]string)
	for _, route := range router.routeSlice {
		if method == route.method && route.regex.MatchString(url) {
			subMatch := route.regex.FindStringSubmatch(url)
			for i, param := range route.params {
				params[param] = subMatch[i+1]
			}
			handler = route.handler
			return params, handler
		}
	}
	return nil, nil
}

// Parse the URL to a regexp and a map of parameters
func (route *route) parseURL(pattern string) (regex *regexp.Regexp, params []string) {
	params = make([]string, 0)
	segments := strings.Split(pattern, "/")
	for i, segment := range segments {
		if strings.HasPrefix(segment, ":") {
			segments[i] = `([\w-%]+)`
			params = append(params, strings.TrimPrefix(segment, ":"))
		}
	}
	regex, _ = regexp.Compile("^" + strings.Join(segments, "/") + "$")
	return regex, params
}
