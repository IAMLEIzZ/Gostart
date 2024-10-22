package gee

import (
	"log"
	"net/http"
	"strings"
)

// 解耦路由
type router struct {
	roots    map[string]*node       // 路由 trie 树
	handlers map[string]HandlerFunc // 路由映射表
}

// newRouter 创建一个路由表
// @return *router
// @author IAMLEIzZ
// @date 2024-10-21 03:02:41
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析路由 允许最后的路径中只包含一个 *
func parsePattern(pattern string) []string {
	paths := make([]string, 0)
	vs := strings.Split(pattern, "/")
	for _, path := range vs {
		if path != "" {
			paths = append(paths, path)
			if path[0] == '*' {
				break
			}
		}
	}

	return paths
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
	// 解析路由
	key := method + "-" + pattern
	paths := parsePattern(pattern)
	// 获取当前请求方式对应的 trie 树的根节点
	_, ok := r.roots[method]

	// 如果当前根节点为 nil，则代表当前树中还没有这个方法对应的路由树，则为这个方法创建一个空的根节点
	if !ok {
		r.roots[method] = &node{}
	}
	// 插入该节点
	r.roots[method].insert(pattern, paths, 0)
	// 插入节点对应方法
	r.handlers[key] = handler
}

//  根据 pattern 获得找到的根节点，以及返回对应解析的参数
func (r *router) getRouter(method string, pattern string) (*node, map[string]string){
	// 解析参数
	paths := parsePattern(pattern)
	param := make(map[string]string)
	// 获取方法对应的树
	root := r.roots[method]
	if root == nil {
		return nil, nil
	}
	// 搜索
	n := root.search(paths, 0)
	// 树中有对应路径
	if n != nil {
		// 树中有对应的路径，n.pattern 即为对应路径模式
		modePath := parsePattern(n.pattern)
		for index, path := range modePath {
			// 如果当前path[0]为 ":" 开头，则往 param 中加入对应参数
			if path[0] == ':' {
				param[path[1:]] = paths[index]
				// TODO: why len(path) > 1 ?
			} else if path[0] == '*' && len(path) > 1{		// 如果当前 path 时一个通配符，则直接把结尾所有路径放入 param 中
				param[path[1:]] = strings.Join(paths[index:], "/")
				break
			}
		}
		return n, param
	}
	
	return nil, nil
}

// handle
// @receiver r
// @param c
// @author IAMLEIzZ
// @date 2024-10-21 03:13:03
func (r *router) handle(c *Context) {
	// 在这里 请求进来后将请求转化为 map 中 key 的格式 {GET-/、GET-/hello、POST-/hello}，然后对 map 进行访问
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
