/*
-*- encoding: utf-8 -*-
@File    :   singleflight.go
@Time    :   2022/10/29 13:20:18
@Author  :   victor2022
@Version :   1.0
@Desc    :   use singleflight to implete batch request

	Singleflight has been impleted in package sync.singleflight,
	In this file, singleflight will be impleted by ourselves.
	The basic structure in this file refers to sync.singleflight.
*/
package singleflight

import "sync"

/*
@Time    :   2022/10/29 13:23:51
@Author  :   victor2022
@Desc    :   代表一个进行中或者已经结束的请求
*/
type call struct {
	wg  sync.WaitGroup // 用于等待一组线程的结束
	val interface{}    // 调用结果
	err error
}

/*
@Time    :   2022/10/29 13:23:34
@Author  :   victor2022
@Desc    :   主结构，管理不同的请求
*/
type Group struct {
	mu sync.Mutex       // 实现map操作时的并发控制
	m  map[string]*call // 请求map，相同的请求指向同一个call
}

/*
@Time    :   2022/10/29 13:32:55
@Author  :   victor2022
@Desc    :   实现了请求和并和调用

	从接收当前请求到响应请求的时间里，0
	不论Do方法被调用了多少次，只会执行一次fn函数，从而实现短时间内多次相同请求的合并
*/
func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	// 初始化map
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 查找对应的call并进行处理
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() // 若已经存在了一个call，说明已经发出请求了，只需要等待结果返回即可
		return c.val, c.err
	}
	// 若不存在对应的call，则创建对应的call
	c := new(call)
	c.wg.Add(1) // 锁加1
	g.m[key] = c
	g.mu.Unlock()
	// 调用函数获取值
	c.val, c.err = fn()
	// 手动结束
	c.wg.Done() // 锁减一
	// 删除key
	g.mu.Lock()
	delete(g.m, key) // 从map中删除对应的call
	g.mu.Unlock()

	return c.val, c.err
}
