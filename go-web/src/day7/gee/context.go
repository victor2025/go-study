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
	// 中间件
	handlers []HandlerFunc
	index int
	// 持有engine，方便页面渲染
	engine *Engine
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
		index: -1,
	}
}

// 执行之后的处理流程
// 中间件之间递归调用
func (c *Context) Next(){
	c.index++
	s := len(c.handlers)
	// 执行接下来的所有中间件
	for ; c.index<s; c.index++{
		// 执行下一个处理流程
		c.handlers[c.index](c)
	}
}

// 获取路径中包含的参数
func (c *Context) Param(key string) string {
	value := c.Params[key]
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

/*
@Time    :   2022/10/23 21:31:05
@Author  :   victor2022 
@Desc    :   设置响应体为html形式，传入模板名称以及元数据
*/
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader(CONTENT_TYPE, "text/html")
	c.Status(code)
	// 使用engine中指定名称的模板处理数据
	if err:= c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err!=nil{
		c.Fail(500, err.Error())
	}
}

func (c *Context) Fail(code int, err string){
	// 跳过之后所有处理器
	c.index = len(c.handlers)
	c.JSON(code, H{
		"message":err,
	})
}
