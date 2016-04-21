package golf

import (
	"fmt"
	"sort"
	"strings"
)

type Node struct {
	text     string
	names    map[string]int
	handler  Handler

	parent   *Node
	colon    *Node

	children nodes
	start    byte
	max      byte
	indices  []uint8
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

	if path == ":" {
		if n.colon == nil {
			n.colon = &Node{text: ":"}
		}
		return n.colon, 0, 0
	}

	for i, node := range n.children {
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
					ccNode := &Node{text: path[0:j], children: nodes{node, &Node{text: path[j:]}}}
					node.text = node.text[j:]
					n.children[i] = ccNode
					return ccNode.children[1], 0, i
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
			n.children = append(n.children, currentNode)
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

func (n *Node) findRoute(urlPath string) (*Node, error) {

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

func (n *Node) optimizeRoutes() {

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

func (n *Node) finalize() {
	if len(n.children) > 0 {
		for i := 0; i < len(n.children); i++ {
			n.children[i].finalize()
		}
	}
	if n.colon != nil {
		n.colon.finalize()
	}
	*n = Node{}
}

func (n *Node) string(col int) string {
	var str = "\n" + strings.Repeat(" ", col) + n.text + " -> "
	col += len(n.text) + 4
	for i := 0; i < len(n.indices); i++ {
		if j := n.indices[i]; j != 0 {
			str += n.children[j-1].string(col)
		}
	}
	if n.colon != nil {
		str += n.colon.string(col)
	}
	return str
}

func (n *Node) String() string {
	if n.text == "" {
		return n.string(0)
	}
	col := len(n.text) + 4
	return n.text + " -> " + n.string(col)
}
