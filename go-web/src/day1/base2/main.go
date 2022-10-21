/*
* -*- encoding: utf-8 -*-
@File    :   main.go
@Time    :   2022/10/19 22:00:41
@Author  :   victor2022
@Version :   1.0
@Desc    :   None
*
*/
package main

import (
	"fmt"
	"net/http"
)

// 所有的请求的统一处理器
type Engine struct{}

// 实现ServeHttp接口，所有的请求都会经过该方法进行处理
// 可以实现日志或者异常处理等统一的行为
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q \n", req.URL.Path)
	case "/hello":
		fmt.Fprintf(w, "Hello World!!!")
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	http.ListenAndServe(":9000", engine)
}
