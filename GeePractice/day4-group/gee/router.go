package gee

import (
	"fmt"
	"net/http"
	"strings"
)

// roots key eg, roots['GET'] roots['POST'];
// handlers key eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book'];
type router struct {
	roots    map[string]*node       // 路由树根节点的集合,一种http方法对应一棵树
	handlers map[string]HandlerFunc // 路由请求到处理函数的映射
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路由
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/") // "/a/b" -> [a, *b]

	parts := make([]string, 0)
	for _, part := range vs {
		if part != "" {
			parts = append(parts, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 1. 解析路径
	parts := parsePattern(pattern)

	// 2. 添加到路径到前缀树
	// 2.1 判断method对应的前缀树的根节点是否存在，不存在则创建
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 2.2 插入前缀树
	r.roots[method].insert(pattern, parts, 0)

	// 3. 更新handlers map
	key := fmt.Sprintf("%s-%s", method, pattern)
	r.handlers[key] = handler
}

// 获取请求路径对应的前缀树节点和路由参数map
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	// 1. 判断方法前缀树是否存在
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 2. 获取path对应的前缀树节点
	// parts := parsePattern(path) 因为请求的path是确定的，直接分割就行
	searchParts := strings.Split(path, "/")
	searchParts = searchParts[1:] // 去掉空串""
	n := root.search(searchParts, 0)

	if n == nil {
		return nil, nil
	}

	// 3. 获取路由参数
	params := make(map[string]string)
	parts := parsePattern(n.pattern)
	// path: /a/b/c/d
	// searchParts: [a, b, c, d]
	// n.pattern: /a/:id/*action
	// parts: [a, :id, *action]
	for index, part := range parts {
		if part[0] == ':' {
			params[part[1:]] = searchParts[index]
		}
		if part[0] == '*' {
			params[part[1:]] = strings.Join(searchParts[index:], "/")
			break
		}
	}

	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n == nil {
		c.String(http.StatusNotFound, "%s页面不存在\n", c.Path)
		return
	}

	c.Params = params
	key := fmt.Sprintf("%s-%s", c.Method, n.pattern)
	r.handlers[key](c)
}
