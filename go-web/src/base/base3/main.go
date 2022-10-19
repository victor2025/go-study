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
	"fmt"
	"net/http"

	"gee"
)

func main() {
	// 创建gee实例
	geeServer := gee.New()
	// 添加路由
	geeServer.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})

	geeServer.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	// 开启服务器
	geeServer.Run(":9000")
}
