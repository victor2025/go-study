/*
-*- encoding: utf-8 -*-
@File    :   main.go
@Time    :   2022/10/19 22:31:44
@Author  :   victor2022
@Version :   1.0
@Desc    :   None
*/
package main

import (
	"gee"
	"net/http"
)

func main() {
	// 创建gee实例
	geeServer := gee.New()
	// 添加路由
	// curl -i http://localhost:9000/
	geeServer.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	// curl "http://localhost:9000/hello?name=geektutu"
	geeServer.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	// curl "http://localhost:9000/login" -X POST -d 'username=geektutu&password=1234'
	geeServer.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	// 开启服务器
	geeServer.Run(":9000")
}
