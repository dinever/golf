package golf

import (
	"sort"
	"strings"
)

type Node struct {
	text    string
	names   map[string]int
	handler Handler

	parent   *Node
	wildcard *Node
	colon    *Node

	nodes   nodes
	start   byte
	max     byte
	indices []uint8
}

type nodes []*Node

func (s nodes) Len() int {
	return len(s)
}

func (s nodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s nodes) Less(i, j int) bool {
	return s[i].text[0] < s[j].text[0]
}

func (n *Node) matchNode(path string) (*Node, int8, int) {

	if path == "*" {
		if n.wildcard == nil {
			n.wildcard = &Node{text: "*"}
		}
		return n.wildcard, 0, 0
	}

	if path == ":" {
		if n.colon == nil {
			n.colon = &Node{text: ":"}
		}
		return n.colon, 0, 0
	}

	for i, node := range n.nodes {
		if node.text[0] == path[0] {

			maxLength := len(node.text)
			pathLength := len(path)
			var pathCompare int8

			if pathLength > maxLength {
				pathCompare = 1
			} else if pathLength < maxLength {
				maxLength = pathLength
				pathCompare = -1
			}

			for j := 0; j < maxLength; j++ {
				if path[j] != node.text[j] {
					ccNode := &Node{text: path[0:j], nodes: nodes{node, &Node{text: path[j:]}}}
					node.text = node.text[j:]
					n.nodes[i] = ccNode
					return ccNode.nodes[1], 0, i
				}
			}

			return node, pathCompare, i
		}
	}

	return nil, 0, 0
}

func (n *Node) addRoute(parts []string, names map[string]int, handler Handler) {

	var (
		tmpNode     *Node
		currentNode *Node
		loop        = true
	)

	currentNode, result, i := n.matchNode(parts[0])

	for loop == true {
		if currentNode == nil {
			currentNode = &Node{text: parts[0]}
			n.nodes = append(n.nodes, currentNode)
		} else if result == 1 {
			//
			parts[0] = parts[0][len(currentNode.text):]
			tmpNode, result, i = currentNode.matchNode(parts[0])
			n = currentNode
			currentNode = tmpNode
			continue
		} else if result == -1 {
			tmpNode := &Node{text: parts[0]}
			currentNode.text = currentNode.text[len(tmpNode.text):]
			tmpNode.nodes = nodes{currentNode}
			n.nodes[i] = tmpNode
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

func (n *Node) findRoute(urlPath string) (*Node, int) {

	urlByte := urlPath[0]
	pathLen := len(urlPath)

	if urlByte >= n.start && urlByte <= n.max {
		if i := n.indices[urlByte-n.start]; i != 0 {
			matched := n.nodes[i-1]
			nodeLen := len(matched.text)
			if nodeLen < pathLen {
				if matched.text == urlPath[:nodeLen] {
					if matched, wildcard := matched.findRoute(urlPath[nodeLen:]); matched != nil {
						return matched, wildcard
					}
				}
			} else if matched.text == urlPath {
				if matched.handler == nil && matched.wildcard != nil {
					return matched.wildcard, 0
				}
				return matched, 0
			}
		}
	}

	if n.colon != nil && pathLen != 0 {
		i := strings.IndexByte(urlPath, '/')
		if i > 0 {
			if cNode, wildcard := n.colon.findRoute(urlPath[i:]); cNode != nil {
				return cNode, wildcard
			}
		} else if n.colon.handler != nil {
			return n.colon, 0
		}
	}

	if n.wildcard != nil {
		return n.wildcard, pathLen
	}

	return nil, 0
}

func (n *Node) optimizeRoutes() {

	if len(n.nodes) > 0 {
		sort.Sort(n.nodes)
		for i := 0; i < len(n.indices); i++ {
			n.indices[i] = 0
		}

		n.start = n.nodes[0].text[0]
		n.max = n.nodes[len(n.nodes)-1].text[0]

		for i := 0; i < len(n.nodes); i++ {
			cNode := n.nodes[i]
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

	if n.wildcard != nil {
		n.wildcard.parent = n
		n.wildcard.optimizeRoutes()
	}
}

func (_node *Node) finalize() {
	if len(_node.nodes) > 0 {
		for i := 0; i < len(_node.nodes); i++ {
			_node.nodes[i].finalize()
		}
	}
	if _node.colon != nil {
		_node.colon.finalize()
	}
	if _node.wildcard != nil {
		_node.wildcard.finalize()
	}
	*_node = Node{}
}

func (_node *Node) string(col int) string {
	var str = "\n" + strings.Repeat(" ", col) + _node.text + " -> "
	col += len(_node.text) + 4
	for i := 0; i < len(_node.indices); i++ {
		if j := _node.indices[i]; j != 0 {
			str += _node.nodes[j-1].string(col)
		}
	}
	if _node.colon != nil {
		str += _node.colon.string(col)
	}
	if _node.wildcard != nil {
		str += _node.wildcard.string(col)
	}
	return str
}

func (_node *Node) String() string {
	if _node.text == "" {
		return _node.string(0)
	}
	col := len(_node.text) + 4
	return _node.text + " -> " + _node.string(col)
}
