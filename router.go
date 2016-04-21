package golf

import (
	"fmt"
	"strings"
)

// Handler is the type of the handler function that Golf accepts.
type Handler func(ctx *Context)

// ErrorHandlerType is the type of the function that handles error in Golf.
type ErrorHandlerType func(ctx *Context, data ...map[string]interface{})

type router struct {
	trees map[string]*Node
}

func newRouter() *router {
	return &router{trees: make(map[string]*Node)}
}

func splitURLpath(path string) (parts []string, names map[string]int) {

	var (
		nameidx      = -1
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

func (router *router) Finalize() {
	for _, _node := range router.trees {
		_node.finalize()
	}
}

func (router *router) FindRoute(method string, path string) (Handler, Parameter, error) {
	node := router.trees[method]
	if node == nil {
		return nil, Parameter{}, fmt.Errorf("Can not find route")
	}
	matchedNode, err := node.findRoute(path)
	if err != nil {
		return nil, Parameter{}, err
	}
	return matchedNode.handler, Parameter{Node: matchedNode, path: path}, err
}

func (router *router) AddRoute(method string, path string, handler Handler) {
	var (
		rootNode *Node
		ok       bool
	)
	parts, names := splitURLpath(path)
	if rootNode, ok = router.trees[method]; !ok {
		rootNode = &Node{}
		router.trees[method] = rootNode
	}
	rootNode.addRoute(parts, names, handler)
	rootNode.optimizeRoutes()
}

func (router *router) String() string {
	var lines string
	for method, _node := range router.trees {
		lines += method + " " + _node.String()
	}
	return lines
}

//Parameter holds the parameters matched in the route
type Parameter struct {
	*Node        // matched node
	path  string // url path given
	cached map[string]string
}

//Len returns number arguments matched in the provided URL
func (p *Parameter) Len() int {
	return len(p.names)
}

//ByName returns the url parameter by name
func (p *Parameter) ByName(name string) (string, error) {
	if i, has := p.names[name]; has {
		return p.findParam(i)
	}
	return "", fmt.Errorf("Parameter not found")
}

//findParam walks up the matched node looking for parameters returns the last parameter
func (p *Parameter) findParam(idx int) (string, error) {
	index := len(p.names) - 1
	urlPath := p.path
	pathLen := len(p.path)
	node := p.Node

	for node != nil {
		if node.text[0] == ':' {
			ctn := strings.LastIndexByte(urlPath, '/')
			if ctn == -1 {
				return "", fmt.Errorf("Parameter not found")
			}
			pathLen = ctn + 1
			if index == idx {
				return urlPath[pathLen:], nil
			}
			index--
		} else {
			pathLen -= len(node.text)
		}
		urlPath = urlPath[0:pathLen]
		node = node.parent
	}
	return "", fmt.Errorf("Parameter not found")
}
