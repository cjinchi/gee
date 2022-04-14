package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc)}
}

func parsePattern(pattern string) []string {
	parts := make([]string, 0)
	for _, item := range strings.Split(pattern, "/") {
		if item == "" {
			continue
		}
		parts = append(parts, item)
		if item[0] == '*' {
			// at most one wildcard
			break
		}
	}
	return parts

}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	pathParts := parsePattern(path)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	n := root.search(pathParts, 0)
	if n == nil {
		return nil, nil
	}

	params := make(map[string]string)
	patternParts := parsePattern(n.pattern)
	for i, part := range patternParts {
		if part[0] == ':' {
			params[part[1:]] = pathParts[i]
		}
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(pathParts[i:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	} else {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	}
}
