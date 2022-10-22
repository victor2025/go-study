/*
-*- encoding: utf-8 -*-
@File    :   router.go
@Time    :   2022/10/21 23:41:41
@Author  :   victor2022
@Version :   1.0
@Desc    :

	router of gee
	use trie store route infomation
*/
package gee

import (
	"log"
	"net/http"
	"strings"
)

// 路由存储结构
type router struct {
	// 存储不同请求方法的根节点，不同路由方法对应不同的trie
	roots map[string]*node
	// 存储不同路由对应的处理方法
	handlers map[string]HandlerFunc
}

// 路由存储结构的工厂方法
func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 解析访问地址，整个地址中只能有一个通配符
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	// 存储parts
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			// 通配符只能放在整个匹配串的结尾
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	// 解析访问地址
	parts := parsePattern(pattern)
	// 生成对应的key
	key := method + "-" + pattern
	// 若当前方法对应的根节点不存在，则创建一个空节点作为根节点
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	// 添加节点
	r.roots[method].insert(pattern, parts, 0)
	// 添加对应的处理器
	r.handlers[key] = handler
}

// 获取访问路径对应的路由节点和参数
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// 解析路径
	searchParts := parsePattern(path)
	// 创建参数map
	params := make(map[string]string)
	// 找到请求方法对应的根节点
	root, ok := r.roots[method]
	// 若对应方法没有根节点，则表明没有创建对应路由
	if !ok {
		return nil, nil
	}

	// 找到对应节点
	n := root.search(searchParts, 0)
	// 若找到了对应的节点
	if n != nil {
		// 找到对应的访问地址
		parts := parsePattern(n.pattern)
		// 遍历所有part，解析参数
		for index, part := range parts {
			// 向参数列表中存放通配路径对应的值
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
	}
	return n, params
}

// 处理请求上下文
func (r *router) handle(c *Context) {
	// 获取当前地址对应的节点
	n, params := r.getRoute(c.Method, c.Path)
	// 若可以查到对应节点，则表明有对应路由
	if n != nil {
		// 向context中设置路径中包含的参数
		c.Params = params
		// 使用对应的handler处理context
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		// 若找不到对应的节点，则表明没有对应的路由
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
