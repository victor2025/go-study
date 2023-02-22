/*
-*- encoding: utf-8 -*-
@File    :   unit.go
@Time    :   2022/10/22 19:15:49
@Author  :   victor2022
@Version :   1.0
@Desc    :   gee的单元测试
*/
package gee

import (
	"fmt"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.addRoute("GET", "/", nil)
	r.addRoute("GET", "/hello/:name", nil)
	r.addRoute("GET", "/hello/b/c", nil)
	r.addRoute("GET", "/hi/:name", nil)
	r.addRoute("GET", "/assets/*filepath", nil)
	return r
}

// 测试解析路径函数
func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

// 测试获取路由
func TestGetRoute(t *testing.T) {
	r := newTestRouter()
	n, ps := r.getRoute("GET", "/hello/geektutu")

	// 测试节点的查找
	if n == nil {
		t.Fatal("nil shouldn't be returned")
	}

	// 测试任意值匹配
	if n.pattern != "/hello/:name" {
		t.Fatal("should match /hello/:name")
	}

	// 测试路径参数的获取
	if ps["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}

	fmt.Printf("matched path: %s, params['name']: %s\n", n.pattern, ps["name"])

}
