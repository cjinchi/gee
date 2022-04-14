package gee

import (
	"strings"
)

type node struct {
	pattern    string
	part       string
	children   []*node
	isWildcard bool
}

func (n *node) matchOneChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWildcard {
			return child
		}
	}
	return nil
}

func (n *node) matchAllChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWildcard {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchOneChild(part)
	if child == nil {
		child = &node{part: part, isWildcard: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchAllChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
