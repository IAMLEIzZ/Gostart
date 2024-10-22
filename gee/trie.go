package gee

import "strings"

type node struct {
	pattern string	// 完整的匹配模式	如/a/b/c
	part string		// 当前节点代表的路径部分值	/a
	children []*node 	// 当前节点的子节点
	isWild bool  	//  是否精准匹配
}

// 精准寻找第一个匹配的孩子节点，用于插入某个节点
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 寻找一组匹配到的节点，用于路由搜索
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild{
			nodes = append(nodes, child)
		}
	}

	return nodes
}

// 插入一个节点
func (n *node) insert(pattern string, paths []string, height int) {
	// 如果深度到达当前 path 的长度，则退出递归
	if len(paths) == height {
		n.pattern = pattern
		return 
	}

	path := paths[height]
	child := n.matchChild(path)

	// 如果没有搜索到当前节点，则直接创建一个新节点插入
	if child == nil {
		// 注意这里是直接给 childnew 了一个新的 node 对象，因为路径插入是沿着一条路到头的，如果不是在原本的基础上 new 的话，下面的
		// child.insert(pattern, paths, height + 1) child 还是 nil，为空，会报错
		child = &node{part: path, isWild: path[0] == ':' || path[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, paths, height + 1)
}

// 路径搜索
func (n *node) search(paths []string, height int) *node {
	// 如果深度到达当前 path 的长度，则退出递归
	if len(paths) == height || strings.HasPrefix(n.part, "*")  {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	path := paths[height]
	children := n.matchChildren(path)

	for _, child := range children {
		result := child.search(paths, height + 1)
		if result != nil {
			return result
		}
	}

	return nil
}