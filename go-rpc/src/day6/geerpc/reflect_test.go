/*
-*- encoding: utf-8 -*-
@File    :   reflect_test.go
@Time    :   2022/11/07 10:46:24
@Author  :   victor2022
@Version :   1.0
@Desc    :   test usage of reflect
*/
package geerpc

import (
	"log"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestReflect(*testing.T) {
	var wg sync.WaitGroup
	// 获取变量的类型
	typ := reflect.TypeOf(&wg)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		argv := make([]string, 0, method.Type.NumIn())
		returns := make([]string, 0, method.Type.NumOut())
		// j 从 1 开始，第 0 个入参是 wg 自己。
		for j := 1; j < method.Type.NumIn(); j++ {
			argv = append(argv, method.Type.In(j).Name())
		}
		for j := 0; j < method.Type.NumOut(); j++ {
			returns = append(returns, method.Type.Out(j).Name())
		}
		log.Printf("func (w *%s) %s(%s) %s",
			typ.Elem().Name(),
			method.Name,
			strings.Join(argv, ","),
			strings.Join(returns, ","))
	}
}
