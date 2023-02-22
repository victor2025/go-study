/*
-*- encoding: utf-8 -*-
@File    :   recovery.go
@Time    :   2022/10/24 10:03:40
@Author  :   victor2022 
@Version :   1.0
@Desc    :   panic recovery of gee
*/
package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

/*
@Time    :   2022/10/24 10:06:31
@Author  :   victor2022 
@Desc    :   recovery handler
*/
func Recovery() HandlerFunc{
	return func(c *Context){
		defer func(){
			if err:=recover(); err!=nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
			
		}()
		c.Next()
	}
}

/*
@Time    :   2022/10/24 10:09:10
@Author  :   victor2022 
@Desc    :   print stack trace for debug
*/
func trace(message string) string{
	// 保存程序计数器
	var pcs [32]uintptr
	// 跳过前三个caller
	n := runtime.Callers(3,pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _,pc := range pcs[:n]{
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d",file, line))
	}
	return str.String()
} 