package Golf

import (
	"regexp"
	"strings"
)

const (
	RouterMethodGet = "GET"
	RouterMethodPost = "POST"
	RouterMethodPut = "PUT"
	RouterMethodDelete = "DELET"
)

type Router struct {
	routeSlice []*Route
}

type Handler func(req Request, res Response)

type Route struct {
	method string
	regex *regexp.Regexp
	params []string
	handler Handler
}

func NewRouter() *Router {
	router := new(Router)
	router.routeSlice = make([]*Route, 0)
	return router
}

func newRoute(method string, pattern string, handler Handler) *Route {
	route := new(Route)
	route.params = make([]string, 0)
	route.regex, route.params = route.parseURL(pattern)
	route.method = method
	route.handler = handler
	return route
}

func (router *Router) Get(pattern string, handler Handler) {
	route := newRoute(RouterMethodGet, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *Router) Post(pattern string, handler Handler) {
	route := newRoute(RouterMethodPost, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *Router) Put(pattern string, handler Handler) {
	route := newRoute(RouterMethodPut, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *Router) Delete(pattern string, handler Handler) {
	route := newRoute(RouterMethodDelete, pattern, handler)
	router.registerRoute(route, handler)
}

func (router *Router) registerRoute(route *Route, handler Handler) {
	router.routeSlice = append(router.routeSlice, route)
}

func (router *Router) match(url string, method string) (params map[string]string, handler Handler) {
	params = make(map[string]string)
	for _, route := range router.routeSlice {
		if method == route.method && route.regex.MatchString(url) {
			subMatch := route.regex.FindStringSubmatch(url)
			for i, param := range route.params {
				params[param] = subMatch[i + 1]
			}
			handler = route.handler
			return
		}
	}
	return nil, nil
}

func (route *Route) parseURL(pattern string) (regex *regexp.Regexp, params []string) {
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
