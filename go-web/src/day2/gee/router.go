/*
-*- encoding: utf-8 -*-
@File    :   router.go
@Time    :   2022/10/21 23:41:41
@Author  :   victor2022
@Version :   1.0
@Desc    :   router of gee
*/
package gee

import (
	"log"
	"net/http"
)

// 路由存储结构
type router struct {
	handlers map[string]HandlerFunc
}

// 路由存储结构的工厂方法
func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 处理请求上下文
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
