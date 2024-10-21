package gee

import (
	"net/http"
)

// 定义一个 handler 方法，是用户定义自己的路由 handler 的统一窗口
type HandlerFunc func(*Context)

// 用户通过获得 engine 对象，往对象中加入自定的路由和处理方式，从而实现静态路由映射
type Engine struct {
	router *router
}

// NewEngine 用户创建一个新的 engine
// @return *Engine
// @author IAMLEIzZ
// @date 2024-10-21 01:11:11
func NewEngine() *Engine {
	return &Engine{router: newRouter()}
}

// addRouter 添加路由信息
// @receiver engine
// @param method
// @param pattern
// @param handler
// @author IAMLEIzZ
// @date 2024-10-21 01:11:23
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRouter(method, pattern, handler)
}

// Get 添加一个 Get 路由
// @receiver engine
// @param pattern
// @param handler
// @author IAMLEIzZ
// @date 2024-10-21 01:13:33
func (engine *Engine) Get(pattern string, handler HandlerFunc) {
	engine.router.addRouter("GET", pattern, handler)
}

// Post 添加一个 Post 路由
// @receiver engine
// @param pattern
// @param handler
// @author IAMLEIzZ
// @date 2024-10-21 01:13:41
func (engine *Engine) Post(pattern string, handler HandlerFunc) {
	engine.router.addRouter("POST", pattern, handler)
}

// Run 启动 http 服务
// @receiver engine
// @param addr
// @return err
// @author IAMLEIzZ
// @date 2024-10-21 01:11:40
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServeHTTP 解析路由, Engine 要实现 http.Handler 接口中的 ServeHTTP 方法
// @receiver engine
// @param w
// @param req
// @author IAMLEIzZ
// @date 2024-10-21 01:11:44
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
