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
	"time"
	"fmt"
	"html/template"
)

type student struct {
	Name string
	Age  int8
}

// 设置渲染函数
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	// 设置渲染函数列表
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	// 加载模板
	r.LoadHTMLGlob("templates/*")
	// 设置静态资源路径
	r.Static("/assets", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}
	// 采用css模板
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	// 采用数组模板
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	// 采用自定义模板
	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})
	// 文件服务器
	r.Static("/file","/home/pi")

	r.Run(":9000")
}
