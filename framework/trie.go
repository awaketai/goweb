package framework

import (
	"fmt"
	"strings"
)

//                  ||
//             /          \
//           user(segment) subject
//           /  \        /     \
//       login logout  name    :id
//                       \       \
//                      age     name(last)
// 1./user/login 2./user/logout 3./subject/name 4./subject/name/age 5/subject/:id/name
type Trie struct {
	root *node // 根节点
}

type node struct {
	isLast   bool                // 代表这个节点是否可以成为最终的路由规则，该节点是否能成为一个独立的uri，是否自身就是一个终极节点
	segment  string              // uri中的字符串，代表这个节点表示的路由中某个段的字符串
	handlers []ControllerHandler // 这个节点中包含的控制器，用于最终加载调用
	childs   []*node             // 子节点
}

// 是否通用开头，即以:开头
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func newNode() *node {
	return &node{
		childs:   make([]*node, 0),
		handlers: make([]ControllerHandler, 0),
	}
}

func NewTree() *Trie {
	return &Trie{
		root: newNode(),
	}
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	// 如果是通配符，则下一层子节点都满足要求
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	// 过滤下一层所有子节点
	for _, cnode := range n.childs {
		if cnode == nil {
			continue
		}
		// if wildcard :
		if isWildSegment(cnode.segment) {
			nodes = append(nodes, cnode)
			// else if there's no wildcard,
			// but segment matched
		} else if cnode.segment == segment {
			nodes = append(nodes, cnode)
		}
	}

	return nodes
}

func (n *node) matchNode(uri string) *node {
	if len(uri) == 0 {
		return nil
	}
	// 0 subject
	// 1 /list/all

	segmentArr := strings.SplitN(uri, "/", 2) // 这样截取出来有空格，可以优化
	segment := segmentArr[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	// match next layer
	cnodes := n.filterChildNodes(segment)
	if len(cnodes) == 0 {
		return nil
	}

	// if has only one segment,it must be the last flag
	if len(segmentArr) == 1 {
		// if the node has isLast flag
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		return nil
	}

	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segmentArr[1])
		if tnMatch != nil {
			return tnMatch
		}
	}

	return nil
}

func (tree *Trie) AddRouter(uri string, handlers []ControllerHandler) error {
	n := tree.root
	if n.matchNode(uri) != nil {
		return fmt.Errorf("route exists:" + uri)
	}

	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := index == len(segments)-1
		var objNode *node // suitable node
		childNodes := n.filterChildNodes(segment)
		// if has child node
		if len(childNodes) > 0 {
			// has same node
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			// create new node
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers
			}
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}

		n = objNode
	}

	return nil
}

func (tree *Trie) FindHandler(uri string) []ControllerHandler {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}

	return matchNode.handlers
}
