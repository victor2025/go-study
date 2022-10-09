package main

import (
	"fmt"
	"time"
)

func PrintNum(name string, c chan string) {
	for i := 0; i < 10; i++ {
		fmt.Printf("%s--%d\n", name, i)
		time.Sleep(500 * time.Millisecond)
	}
	c <- name
}

func main() {
	// 创建信道(带有类型的管道)进行同步，可以设置信道长度
	c := make(chan string)
	// 创建两个routine进行操作
	go PrintNum("routine-1", c)
	go PrintNum("routine-2", c)
	fmt.Printf("%s--%s", <-c, <-c)
}
