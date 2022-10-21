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
	"fmt"
	"log"
	"net/http"
)

// 定义HandlerFunc
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 实现ServeHTTP接口，保存路由信息
type Engine struct {
	router map[string]HandlerFunc
}

// 创建Engine的方法，工厂模式
func New() *Engine {
	// 返回生成的Engine
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 添加路由的方法，属于Engine
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	// 由请求方法和匹配字符生成路由key
	key := method + "-" + pattern
	// 将当前处理器放入路由中
	engine.router[key] = handler
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
	// 生成当前方法当前路径的key
	key := req.Method + "-" + req.URL.Path
	// 查找对应的处理器
	if handler, ok := engine.router[key]; ok {
		// 存在对应的处理器，则进行处理
		handler(w, req)
	} else {
		// 不存在对应的处理器，则打印404错误
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
