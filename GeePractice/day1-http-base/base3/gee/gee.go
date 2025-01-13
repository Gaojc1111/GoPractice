package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc 将
type HandlerFunc func(w http.ResponseWriter, req *http.Request)

//func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	f(w, req)
//}

type Engine struct {
	router map[string]HandlerFunc // 路由到handler的映射
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := fmt.Sprintf("%s-%s", method, pattern)
	//key := method + "-" + pattern
	engine.router[key] = handler
}

func (engine *Engine) Get(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) Post(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP
// 疑惑：为什么是engine去实现ServeHTTP，而不是像net/http包一样，通过HandlerFunc去实现？
// http.ListenAndServe(":9999", engine)
// 个人理解：
// 1.首先，ServeHTTP是Handler接口包含的方法：用来处理路由，接收HTTP请求，返回HTTP响应；
// 2.net/http包提供的路由匹配，比如HandleFunc, 都是一个pattern（如："/hello"）对应一个处理函数(如：helloHandler（w ..., r ...）)
// 3.所以每个处理函数都得去实现ServeHTTP, 因此net/http包提供了HandlerFunc适配器
// 4.现在我们自己定义的engine，通过map聚合了所有的处理函数，所有只需要engine去实现ServeHTTP， 然后在ServeHTTP中去进行路由匹配就行了。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := fmt.Sprintf("%s-%s", req.Method, req.URL.Path)
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
