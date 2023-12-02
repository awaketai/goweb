package framework2

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root: root}
}

type node struct {
	isLast  bool // if the node become the last match;or the node can be a independent url;if the last node
	segment string
	handlers []ControllerHandler // the node controller
	childs  []*node
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, "")
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	if isWildSegment(segment) {
		return n.childs
	}
	nodes := make([]*node, 0, len(n.childs))
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			// have wildcard
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			// not the whildcard,match the url
			nodes = append(nodes, cnode)
		}
	}

	return nodes
}

func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	// next node
	cnodes := n.filterChildNodes(segment)
	if len(cnodes) == 0 {
		return nil
	}
	// only one segment
	if len(segment) == 1 {
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
	}

	// not finally node
	return nil
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
}

// /book/list
// /book/:id conflict
// /book/:id/name
// /book/:student/age
// /:/user/name
// /:user/name/:age conflict
func (t *Tree) AddRouter(uri string, handler []ControllerHandler) error {
	n := t.root
	if n.matchNode(uri) != nil {
		return errors.New("route exists:" + uri)
	}
	segments := strings.Split(uri, "/")
	for i, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		isLast := i == len(segments)-1
		// if have suitable node
		var objNode *node
		childNodes := n.filterChildNodes(segment)
		if len(childNodes) > 0 {
			// have same child node,select
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}
		if objNode == nil {
			// creat node
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handlers = handler
			}
		}
		n = objNode
	}

	return nil
}

func (t *Tree) FindHandler(uri string) []ControllerHandler {
	matchNode := t.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}

	return matchNode.handlers
}
