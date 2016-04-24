package golf

import (
	"fmt"
	"sort"
	"strings"
)

type node struct {
	text    string
	names   map[string]int
	handler HandlerFunc

	parent *node
	colon  *node

	children nodes
	start    byte
	max      byte
	indices  []uint8
}

type nodes []*node

func (s nodes) Len() int {
	return len(s)
}

func (s nodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s nodes) Less(i, j int) bool {
	return s[i].text[0] < s[j].text[0]
}

func (n *node) matchNode(path string) (*node, int8, int) {

	if path == ":" {
		if n.colon == nil {
			n.colon = &node{text: ":"}
		}
		return n.colon, 0, 0
	}

	for i, child := range n.children {
		if child.text[0] == path[0] {

			maxLength := len(child.text)
			pathLength := len(path)
			var pathCompare int8

			if pathLength > maxLength {
				pathCompare = 1
			} else if pathLength < maxLength {
				maxLength = pathLength
				pathCompare = -1
			}

			for j := 0; j < maxLength; j++ {
				if path[j] != child.text[j] {
					ccNode := &node{text: path[0:j], children: nodes{child, &node{text: path[j:]}}}
					child.text = child.text[j:]
					n.children[i] = ccNode
					return ccNode.children[1], 0, i
				}
			}

			return child, pathCompare, i
		}
	}

	return nil, 0, 0
}

func (n *node) addRoute(parts []string, names map[string]int, handler HandlerFunc) {

	var (
		tmpNode     *node
		currentNode *node
	)

	currentNode, result, i := n.matchNode(parts[0])

	for {
		if currentNode == nil {
			currentNode = &node{text: parts[0]}
			n.children = append(n.children, currentNode)
		} else if result == 1 {
			parts[0] = parts[0][len(currentNode.text):]
			tmpNode, result, i = currentNode.matchNode(parts[0])
			n = currentNode
			currentNode = tmpNode
			continue
		} else if result == -1 {
			tmpNode := &node{text: parts[0]}
			currentNode.text = currentNode.text[len(tmpNode.text):]
			tmpNode.children = nodes{currentNode}
			n.children[i] = tmpNode
			currentNode = tmpNode
		}
		break
	}

	if len(parts) == 1 {
		currentNode.handler = handler
		currentNode.names = names
		return
	}

	currentNode.addRoute(parts[1:], names, handler)
}

func (n *node) findRoute(urlPath string) (*node, error) {

	urlByte := urlPath[0]
	pathLen := len(urlPath)

	if urlByte >= n.start && urlByte <= n.max {
		if i := n.indices[urlByte-n.start]; i != 0 {
			matched := n.children[i-1]
			nodeLen := len(matched.text)
			if nodeLen < pathLen {
				if matched.text == urlPath[:nodeLen] {
					if matched, _ := matched.findRoute(urlPath[nodeLen:]); matched != nil {
						return matched, nil
					}
				}
			} else if matched.text == urlPath {
				return matched, nil
			}
		}
	}

	if n.colon != nil && pathLen != 0 {
		i := strings.IndexByte(urlPath, '/')
		if i > 0 {
			if cNode, err := n.colon.findRoute(urlPath[i:]); cNode != nil {
				return cNode, err
			}
		} else if n.colon.handler != nil {
			return n.colon, nil
		}
	}

	return nil, fmt.Errorf("Can not find route")
}

func (n *node) optimizeRoutes() {

	if len(n.children) > 0 {
		sort.Sort(n.children)
		for i := 0; i < len(n.indices); i++ {
			n.indices[i] = 0
		}

		n.start = n.children[0].text[0]
		n.max = n.children[len(n.children)-1].text[0]

		for i := 0; i < len(n.children); i++ {
			cNode := n.children[i]
			cNode.parent = n

			cByte := int(cNode.text[0] - n.start)
			if cByte >= len(n.indices) {
				n.indices = append(n.indices, make([]uint8, cByte+1-len(n.indices))...)
			}
			n.indices[cByte] = uint8(i + 1)
			cNode.optimizeRoutes()
		}
	}

	if n.colon != nil {
		n.colon.parent = n
		n.colon.optimizeRoutes()
	}
}
