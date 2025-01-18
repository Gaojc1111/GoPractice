package gee

import "strings"

type node struct {
	pattern  string  // 完整路由 /p/:lang/doc or tutorial or intro
	part     string  // 路由的一部分 :lang
	children []*node // 该部分的子节点: [doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配， part含有 ":" or "*" 时为 true
}

// pattern: /p/:name/ttt
// parts: [p, :name, ttt]
// height: 当前处理的是第几个字符串
func (n *node) insert(pattern string, parts []string, height int) {
	// 因为height是从0计数，相等时，说明所有的part处理完了。
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	// 1. 获取当前part
	part := parts[height]

	// 2. 判断当前节点的子节点是否包含part，没有则创建
	child := n.matchChild(part)
	if child == nil { // 两种情况，1：children是空的 2：children中不包含
		child = &node{
			pattern:  "", // 只有路由结尾才存储完整路径
			part:     part,
			children: nil, // slice 零值可用
			isWild:   part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child) // 因为slice零值可用，即使是nil也可以append，append内部会特殊处理
	}

	// 3. 递归子节点
	child.insert(pattern, parts, height+1)
}

// parts: [a, b, c, d]
func (n *node) search(parts []string, height int) *node {
	// 如果已经全部匹配完，或者有*
	// 要使用strings.HasPrefix(n.part, "*")而非n.part[0] == '*'， 因为part初始化为空串，index超出 index out of range [0] with length 0
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// todo：感觉下面没必要
		//if n.pattern == "" {
		//	return nil
		//}
		return n
	}

	// 1. 判断当前节点是否有子节点能匹配part
	part := parts[height]
	children := n.matchChildren(part)

	// 2. 递归子节点
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// 插入pattern， 精确匹配
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

// 根据path查询前缀树节点，优先精确，后模糊
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes
}
