/*
-*- encoding: utf-8 -*-
@File    :   main.go
@Time    :   2022/11/03 22:32:25
@Author  :   victor2022
@Version :   1.0
@Desc    :   simple client of geerpc
*/
package main

import (
	"fmt"
	"geerpc"
	"geerpc/codec"
	"log"
	"net"
	"sync"
	"time"
)

/*
@Time    :   2022/11/05 23:22:26
@Author  :   victor2022
@Desc    :   开启服务
*/
func startServer(addr chan string) {
	// 监听一个端口
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", lis.Addr())
	addr <- lis.Addr().String()
	geerpc.Accept(lis)
}

func main() {
	log.SetFlags(0)
	// 开启服务提供者
	addr := make(chan string)
	go startServer(addr)
	// 创建服务消费者
	realAddr := <-addr
	// 使用默认配置(gob)
	clientGob, _ := geerpc.Dial("tcp", realAddr)
	defer func() { _ = clientGob.Close() }()
	// 使用json编解码
	clientJson, _ := geerpc.Dial("tcp", realAddr, &geerpc.Option{MagicNumber: 0, CodecType: codec.JsonType})
	defer func() { _ = clientJson.Close() }()

	time.Sleep(time.Second)

	// 开始发送请求
	w1 := make(chan int)
	go call(clientGob, w1)
	w2 := make(chan int)
	go call(clientJson, w2)
	// 等待结束
	<-w1
	<-w2
}

/*
@Time    :   2022/11/05 23:22:15
@Author  :   victor2022
@Desc    :   发起远程调用
*/
func call(client *geerpc.Client, w chan int) {
	now := time.Now().UnixMilli()
	fmt.Printf("now: %v\n", now)
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := client.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
	// 显示当前时间
	dur := time.Now().UnixMilli() - now
	fmt.Printf("dur: %v\n", dur)
	// 提醒结束
	w <- 0
}
