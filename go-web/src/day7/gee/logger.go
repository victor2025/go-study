/*
-*- encoding: utf-8 -*-
@File    :   logger.go
@Time    :   2022/10/23 13:41:49
@Author  :   victor2022 
@Version :   1.0
@Desc    :   logger middleware of gee
*/
package gee

import (
	"time"
	"log"
)

// 创建log函数作为中间件
func Logger() HandlerFunc{
	return func(c *Context){
		// 开始计时
		t:= time.Now()
		// 处理请求
		c.Next()
		// 计算请求时间
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}