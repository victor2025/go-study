/*
-*- encoding: utf-8 -*-
@File    :   trie.go
@Time    :   2022/10/22 18:15:11
@Author  :   victor2022
@Version :   1.0
@Desc    :   trie for dynamic route
*/
package gee

import (
	"strings"
)

// trie的节点结构
type node struct {
	pattern  string  // 当前节点对应路由，只有叶子节点当前属性才不为空
	part     string  // 路由中的部分
	children []*node // 子节点
	isWild   bool    // 标志当前节点是否需要精确匹配，part含有:或者*时为true
}

// 找到匹配当前路由的第一个子节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		// 当前子节点对应部分与路由部分相同
		// 或者子节点支持模糊匹配，则返回第一个匹配成功的子节点
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 找到所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归方式插入新的节点
func (n *node) insert(pattern string, parts []string, height int) {
	// 层数到达了路由的最后一层，则将匹配串修改为当前匹配串，表示当前节点为叶子节点
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	// 没有到达最后一层，则找到并插入新的子节点
	// 找到当前层对应的part
	part := parts[height]
	// 找到当前节点第一个匹配的子节点
	child := n.matchChild(part)
	// 若没有对应的子节点，则插入新的子节点
	if child == nil {
		// 创建新的子节点
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 将新的节点添加到当前节点的子节点中
		n.children = append(n.children, child)
	}
	// 递归地向子节点中插入后续的节点
	child.insert(pattern, parts, height+1)
}

// 递归查找当前路由对应的节点
func (n *node) search(parts []string, height int) *node {
	// 若当前道道最后一层，或者当前层为通配节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// 判断当前节点是否为叶子节点，若不是，则表示匹配失败
		if n.pattern == "" {
			return nil
		}
		// 若是叶子节点，则匹配成功
		return n
	}
	// 找到对应的part
	part := parts[height]
	// 找到所有对应的子节点
	children := n.matchChildren(part)
	// 递归查找对应节点
	for _, child := range children {
		result := child.search(parts, height+1)
		// 若子节点能够成功匹配，则返回结果
		if result != nil {
			return result
		}
	}
	return nil
}
