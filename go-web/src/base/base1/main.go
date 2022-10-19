/*
* -*- encoding: utf-8 -*-
@File    :   main.go
@Time    :   2022/10/19 16:39:42
@Author  :   victor2022
@Version :   1.0
@Desc    :   day1 for gee
*
*/
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// 映射地址并配置对应的处理器
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/hello", helloHandler)
	// 监听端口号
	log.Fatal(http.ListenAndServe(":9000", nil))
}

// 创建处理器函数
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q \n", r.URL.Path)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Fprintf(w, "Header[%q] = %q \n", k, v)
	}
}
