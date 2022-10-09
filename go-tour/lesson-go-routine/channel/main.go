package main

import (
	"fmt"
	"time"
)

func Produce(c chan int) {
	for i := 0; i < 10; i++ {
		c <- i
		time.Sleep(100 * time.Millisecond)
	}
	// 完成后需要关闭，避免主线程一直等待
	close(c)
}

func main() {
	// 创建信道(带有类型的管道)进行同步，可以设置信道长度
	c := make(chan int)
	// 开始放入数据
	go Produce(c)
	// 创建两个routine进行操作
	for i := range c {
		fmt.Println(i)
	}
}
