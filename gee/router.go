package gee

import (
	"log"
	"net/http"
)

// 解耦路由
type router struct {
	handlers map[string]HandlerFunc	// 路由映射表
}

// newRouter 创建一个路由表
// @return *router  
// @author IAMLEIzZ   
// @date 2024-10-21 03:02:41 
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// addRouter 添加一个静态路由映射
// @receiver r 
// @param method 
// @param pattern 
// @param handler 
// @author IAMLEIzZ   
// @date 2024-10-21 03:02:19 
func (r *router) addRouter(method string, pattern string, handler HandlerFunc) {
	// 日志输出当前路由表中的映射关系
	log.Printf("Route %4s - %s", method, pattern)
	path := method + "-" + pattern
	r.handlers[path] = handler
}

// handle 
// @receiver r 
// @param c 
// @author IAMLEIzZ   
// @date 2024-10-21 03:13:03 
func (r *router) handle(c *Context) {
	// 在这里 请求进来后将请求转化为 map 中 key 的格式 {GET-/、GET-/hello、POST-/hello}，然后对 map 进行访问
	key := c.Method + "-" + c.Path

	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}