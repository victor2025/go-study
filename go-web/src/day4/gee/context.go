/*
-*- encoding: utf-8 -*-
@File    :   context.go
@Time    :   2022/10/21 23:12:34
@Author  :   victor2022
@Version :   1.0
@Desc    :   Context of Request and Response
*/
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// map的别名
type H map[string]interface{}

// 上下文类
type Context struct {
	// 原始数据
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	Params map[string]string // 请求参数
	// 返回信息
	StatusCode int
}

// 常量
const CONTENT_TYPE = "Content-Type"

// Context工厂方法
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	// 返回创建的新的上下文
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 获取路径中包含的参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 查找请求表单中的指定值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 解析请求URL中的参数
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置返回状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置响应体为文本形式
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader(CONTENT_TYPE, "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 设置响应体为json字符串形式
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader(CONTENT_TYPE, "application/json")
	c.Status(code)
	// 由输出流创建json编码器
	encoder := json.NewEncoder(c.Writer)
	// 向输出流写入数据
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 设置响应体为二进制数据形式
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// 设置响应体为html形式
func (c *Context) HTML(code int, html string) {
	c.SetHeader(CONTENT_TYPE, "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
