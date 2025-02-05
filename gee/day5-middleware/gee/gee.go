package gee

import (
	"log"
	"net/http"
	"strings"
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
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} // 把自己这个顶层分组加入所有分组的切片
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: group.engine,
		// 每个分组只维护自己的middleware, 因此不需要继承
		//middlewares: group.middlewares,
	}
	group.engine.groups = append(group.engine.groups, newGroup) // engine这个顶层group，存储了所有的分组信息
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

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		// 判断这个路由属于哪个分组
		// 假设一个请求是 /a/b/c, 如果有分组/a, 并且/a有子分组,那么子分组的prefix就是/a/b, 所以会遍历所有的middleware
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
