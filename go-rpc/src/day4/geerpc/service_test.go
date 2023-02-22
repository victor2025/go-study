/*
-*- encoding: utf-8 -*-
@File    :   service_test.go
@Time    :   2022/11/07 15:02:17
@Author  :   victor2022
@Version :   1.0
@Desc    :   tests for invocation of methods in service
*/
package geerpc

import (
	"fmt"
	"reflect"
	"testing"
)

/*
@Time    :   2022/11/07 15:03:30
@Author  :   victor2022
@Desc    :   执行者结构体
*/
type Foo int

/*
@Time    :   2022/11/07 15:03:41
@Author  :   victor2022
@Desc    :   参数结构体
*/
type Args struct {
	Num1, Num2 int
}

// 暴露的方法
func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

// 未暴露的方法
func (f Foo) sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func _assert(condition bool, msg string, v ...interface{}) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

/*
@Time    :   2022/11/07 15:06:48
@Author  :   victor2022
@Desc    :   测试创建新服务
*/
func TestNewService(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	// 判断是否按照要求进行了注册
	_assert(len(s.method) == 1, "wrong service Method, expect 1, but got %d", len(s.method))
	mType := s.method["Sum"]
	_assert(mType != nil, "wrong Method, Sum shouldn't nil")
}

/*
@Time    :   2022/11/07 15:15:44
@Author  :   victor2022
@Desc    :   测试服务调用
*/
func TestMethodType_Call(t *testing.T) {
	var foo Foo
	s := newService(&foo)
	mType := s.method["Sum"]

	// 获取参数的反射对象
	argv := mType.newArgv()
	replyv := mType.newReplyv()
	// 设置输入参数的值
	argv.Set(reflect.ValueOf(Args{Num1: 1, Num2: 3}))
	// 使用service进行调用
	err := s.call(mType, argv, replyv)
	// 判断是否满足要求
	_assert(err == nil && *replyv.Interface().(*int) == 4 && mType.numCalls == 1,
		"failed to call Foo.Sum")
}
