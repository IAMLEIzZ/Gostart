package gee

import (
	"log"
	"net/http"
)

// 定义一个 handler 方法，是用户定义自己的路由 handler 的统一窗口
type HandlerFunc func(*Context)

// 一个分组对象
type RouterGroup struct {
	prefix      string        // 前缀
	middlewares []HandlerFunc // support middleware
	parent      *RouterGroup  // support nesting
	engine      *Engine       // 这样分组对象可以直接调用 engine 的功能
}

// 用户通过获得 engine 对象，往对象中加入自定的路由和处理方式，从而实现静态路由映射
type Engine struct {
	*RouterGroup // 嵌套了 RouterGroup，所以 Engine 自身也可以作为一个路由组， 这样 Engine 可以直接调用 RouterGroup 的私有方法
	router       *router
	groups       []*RouterGroup
}

// NewEngine 用户创建一个新的 engine
// @return *Engine
// @author IAMLEIzZ
// @date 2024-10-21 01:11:11
// 在分组情况下，一个新的 engine 代表着最高权限，后面创建的分组都是在这个 engine 之下的
func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// NewGroup 	创建一个新的分组，一半来说由 engine 对象调用，作为其下的子分组
// @receiver group
// @param prefix
// @return *RouterGroup
// @author IAMLEIzZ
// @date 2024-10-22 05:07:33
func (group *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

// addRouter 添加路由信息
// @receiver engine
// @param method
// @param pattern
// @param handler
// @author IAMLEIzZ
// @date 2024-10-21 01:11:23
func (group *RouterGroup) addRouter(method string, prefix string, handler HandlerFunc) {
	pattern := group.prefix + prefix
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRouter(method, pattern, handler)
}

// Get 添加一个 Get 路由
// @receiver engine
// @param pattern
// @param handler
// @author IAMLEIzZ
// @date 2024-10-21 01:13:33
func (group *RouterGroup) Get(pattern string, handler HandlerFunc) {
	group.addRouter("GET", pattern, handler)
}

// Post 添加一个 Post 路由
// @receiver engine
// @param pattern
// @param handler
// @author IAMLEIzZ
// @date 2024-10-21 01:13:41
func (group *RouterGroup) Post(pattern string, handler HandlerFunc) {
	group.addRouter("POST", pattern, handler)
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
