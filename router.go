package golf

import (
	"fmt"
)

// HandlerFunc is the type of the handler function that Golf accepts.
type HandlerFunc func(ctx *Context)

// ErrorHandlerFunc is the type of the function that handles error in Golf.
type ErrorHandlerFunc func(ctx *Context, data ...map[string]interface{})

type router struct {
	trees map[string]*node
}

func newRouter() *router {
	return &router{trees: make(map[string]*node)}
}

func splitURLPath(path string) (parts []string, names map[string]int) {

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
					panic(fmt.Errorf("Invalid parameter : or * should always be after / - %q", path))
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

func (router *router) FindRoute(method string, path string) (HandlerFunc, Parameter, error) {
	node := router.trees[method]
	if node == nil {
		return nil, Parameter{}, fmt.Errorf("Can not find route")
	}
	matchedNode, err := node.findRoute(path)
	if err != nil {
		return nil, Parameter{}, err
	}
	return matchedNode.handler, Parameter{node: matchedNode, path: path}, err
}

func (router *router) AddRoute(method string, path string, handler HandlerFunc) {
	var (
		rootNode *node
		ok       bool
	)
	if rootNode, ok = router.trees[method]; !ok {
		rootNode = &node{}
		router.trees[method] = rootNode
	}

	parts, names := splitURLPath(path)
	rootNode.addRoute(parts, names, handler)

	if path == "/" {
	} else if path[len(path) - 1] != '/' {
		parts, names := splitURLPath(path + "/")
		rootNode.addRoute(parts, names, handler)
	} else {
		parts, names := splitURLPath(path[:len(path) - 1])
		rootNode.addRoute(parts, names, handler)
	}
	rootNode.optimizeRoutes()
}

//Parameter holds the parameters matched in the route
type Parameter struct {
	*node         // matched node
	path   string // url path given
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

func lastIndexByte(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}

//findParam walks up the matched node looking for parameters returns the last parameter
func (p *Parameter) findParam(idx int) (string, error) {
	index := len(p.names) - 1
	urlPath := p.path
	pathLen := len(p.path)
	node := p.node

	for node != nil {
		if node.text[0] == ':' {
			ctn := lastIndexByte(urlPath, '/')
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
