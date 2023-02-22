/*
-*- encoding: utf-8 -*-
@File    :   gee.go
@Time    :   2022/10/19 22:30:32
@Author  :   victor2022
@Version :   1.0
@Desc    :   None
*/
package gee

import (
	"log"
	"net/http"
)

// 定义HandlerFunc
type HandlerFunc func(*Context)

// 实现ServeHTTP接口，保存路由信息
type Engine struct {
	router *router
}

// 创建Engine的方法，工厂模式
func New() *Engine {
	// 返回生成的Engine
	return &Engine{router: newRouter()}
}

// 添加路由的方法，属于Engine
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// 添加Get请求的方法
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	// 调用公用方法添加路由
	engine.addRoute("GET", pattern, handler)
}

// 添加POST请求的方法
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
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
