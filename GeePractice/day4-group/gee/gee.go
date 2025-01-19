package gee

import (
	"log"
	"net/http"
)

// HandlerFunc 将
type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup                // 组合，实际上只拥有方法
	router       *router        // 路由到handler的映射
	groups       []*RouterGroup // 所有的根分组, 其实就是一个分组单链表切片
}

// RouterGroup
// Engine 和 RouterGroup的关系：
// 每一个RouterGroup包含了分组后的信息，并且可以通过engine字段，去创建路由映射；
// Engine通过组合RouterGroup， 也具备了分组信息，相当于最顶层分组。
type RouterGroup struct {
	prefix      string
	parent      *RouterGroup
	middlewares []HandlerFunc
	engine      *Engine
}

func New() *Engine {
	engine := &Engine{
		router: newRouter(),
		groups: make([]*RouterGroup, 0),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
	}
	group.engine.groups = append(group.engine.groups, newGroup)
	return newGroup
}

// comp 这名字不太理解，比较的意思么？
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %-4s - %s", method, pattern) // "-" 左对齐
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

//func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
//	log.Printf("Route %-4s - %s", method, pattern) // "-" 左对齐
//	engine.router.addRoute(method, pattern, handler)
//}

//func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//	engine.addRoute("GET", pattern, handler)
//}

//func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//	engine.addRoute("POST", pattern, handler)
//}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
