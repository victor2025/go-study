package main

import (
	"fmt"
	"sync"
	"time"
)

// 公共接口
type Counter interface {
	Inc(key string) int
	Value(key string) int
}

// 自带互斥量的计数器
type SafeCounter struct {
	val map[string]int
	mut sync.Mutex
}

// 接口的实现
// 计数自增操作
func (c *SafeCounter) Inc(key string) int {
	// 加锁
	c.mut.Lock()
	// 自增
	c.val[key]++
	// 解锁
	defer c.mut.Unlock()
	return c.val[key]
}

// 获取值操作
func (c *SafeCounter) Value(key string) int {
	c.mut.Lock()
	// 推迟到暂存返回值后执行
	defer c.mut.Unlock()
	return c.val[key]
}

func main() {
	var counter Counter = &SafeCounter{
		val: make(map[string]int),
		mut: sync.Mutex{},
	}
	// 开始测试
	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println(counter.Inc("aa"))
		}
	}()
	go func() {
		for i := 0; i < 100; i++ {
			fmt.Println(counter.Inc("aa"))
		}
	}()
	// 等待程序执行完成
	time.Sleep(2 * time.Second)
	fmt.Println("final value:", counter.Value("aa"))
}
