package golf

import (
	"fmt"
	"net/http"
	"golang.org/x/net/context"
)

type Handler func(c context.Context, w http.ResponseWriter, r *http.Request)

type Router struct {
	Trees map[string]*Node
}

func NewRouter() *Router {
	return &Router{Trees: make(map[string]*Node)}
}

func splitURLpath(path string) (parts []string, names map[string]int) {

	var (
		nameidx      int = -1
		partidx      int
		paramCounter int
	)

	for i := 0; i < len(path); i++ {

		if names == nil {
			names = make(map[string]int)
		}
		// recording name
		if nameidx != -1 {
			//found /
			if path[i] == '/' {
				names[path[nameidx:i]] = paramCounter
				paramCounter++

				nameidx = -1 // switch to normal recording
				partidx = i
			}
		} else {
			if path[i] == ':' || path[i] == '*' {
				if path[i-1] != '/' {
					panic(fmt.Errorf("Inválid parameter : or * comes anwais after / - %q", path))
				}
				nameidx = i + 1
				if partidx != i {
					parts = append(parts, path[partidx:i])
				}
				parts = append(parts, path[i:nameidx])
			}
		}
	}

	if nameidx != -1 {
		names[path[nameidx:]] = paramCounter
		paramCounter++
	} else if partidx < len(path) {
		parts = append(parts, path[partidx:])
	}
	return
}

func (router *Router) Finalize() {
	for _, _node := range router.Trees {
		_node.finalize()
	}
}

func (router *Router) FindRoute(method string, path string) (Handler, Parameter) {
	node := router.Trees[method]
	if node == nil {
		return nil, Parameter{}
	}
	matchedNode, wildcard := node.findRoute(path)
	if matchedNode != nil {
		return matchedNode.handler, Parameter{Node: matchedNode, path: path, wildcard: wildcard}
	}
	return nil, Parameter{}
}

func (router *Router) AddRoute(method string, path string, handler Handler) {
	var (
		rootNode *Node
		ok       bool
	)
	parts, names := splitURLpath(path)
	if rootNode, ok = router.Trees[method]; !ok {
		rootNode = &Node{}
		router.Trees[method] = rootNode
	}
	rootNode.addRoute(parts, names, handler)
	rootNode.optimizeRoutes()
}

func (router *Router) String() string {
	var lines string
	for method, _node := range router.Trees {
		lines += method + " " + _node.String()
	}
	return lines
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, variables := router.FindRoute(r.Method, r.URL.Path)

	if handler != nil {
		handler(w, r, variables)
	} else {
		http.NotFound(w, r)
	}
}
