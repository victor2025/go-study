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
	"encoding/json"
	"fmt"
	"geerpc"
	"geerpc/codec"
	"log"
	"net"
	"time"
)

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
	addr := make(chan string)
	go startServer(addr)

	// 创建简单的rpc客户端
	conn, _ := net.Dial("tcp", <-addr)
	defer func() { _ = conn.Close() }()

	time.Sleep(time.Second)
	// 发送设置
	_ = json.NewEncoder(conn).Encode(geerpc.DefaultOption)
	cc := codec.NewGobCodec(conn)
	// 发送请求并接收响应
	for i := 0; i < 5; i++ {
		h := &codec.Header{
			ServiceMethod: "Foo.Sum",
			Seq:           uint64(i),
		}
		_ = cc.Write(h, fmt.Sprintf("geerpc req %d", h.Seq))
		_ = cc.ReadHeader(h)
		var reply string
		_ = cc.ReadBody(&reply)
		log.Println("reply", reply)
	}
}
