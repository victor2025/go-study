/*
-*- encoding: utf-8 -*-
@File    :   gee.go
@Time    :   2022/10/19 22:30:32
@Author  :   victor2022
@Version :   1.0
@Desc    :   main structs of gee
*/
package gee

import (
	"log"
	"net/http"
)

type (
	// 定义HandlerFunc
	HandlerFunc func(*Context)

	// 定义Group结构，实现分组功能
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // 中间件扩展
		parent      *RouterGroup  // 父Group
		engine      *Engine       // 所有的Group持有一个Engine实例
	}

	// 实现ServeHTTP接口，保存路由信息
	// 将Engine抽象为一个顶层分组
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup // store all groups
	}
)

// 创建Engine的方法，工厂模式
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	// 初始元素只有当前主Group
	engine.groups = []*RouterGroup{engine.RouterGroup}
	// 返回生成的Engine
	return engine
}

// 添加分组的方法
// 所有的分组都共享一个engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		// 拼接当前前缀
		prefix: group.prefix + prefix,
		// 父分组
		parent: group,
		// 持有的engine
		engine: engine,
	}
	// 向group组中添加当前group
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由的方法，属于Group
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// 拼接pattern
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 添加Get请求的方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	// 调用公用方法添加路由
	group.addRoute("GET", pattern, handler)
}

// 添加POST请求的方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 开启服务器
func (engine *Engine) Run(addr string) (err error) {
	// 开始监听端口，提供服务
	res := http.ListenAndServe(addr, engine)
	// 在发生错误时输出日志
	log.Fatal(res)
	// 返回可能的错误
	return res
}

// Engine要实现的接口中的方法
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 创建对应的上下文
	c := newContext(w, req)
	engine.router.handle(c)
}
