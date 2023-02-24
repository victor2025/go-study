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
	"context"
	"geerpc"
	"geerpc/registry"
	"geerpc/xclient"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

type Foo int

type Args struct{ Num1, Num2 int }

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1)) // 睡眠，用来造成调用失败
	*reply = args.Num1 + args.Num2
	return nil
}

// 开启注册中心
func startRegistry(wg *sync.WaitGroup) {
	l, _ := net.Listen("tcp", ":9000")
	// 开启默认注册中心
	registry.HandleHTTP()
	wg.Done()
	_ = http.Serve(l, nil)
}

/*
@Time    :   2022/11/07 15:45:28
@Author  :   victor2022
@Desc    :   开启rpc服务器
*/
func startServer(registryAddr string, wg *sync.WaitGroup) {
	// 注册服务
	var foo Foo
	l, _ := net.Listen("tcp", ":0")
	server := geerpc.NewServer()
	_ = server.Register(&foo)
	// 开始心跳,5s一次
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), time.Duration(5)*time.Second)
	wg.Done()
	server.Accept(l)
}

func foo(xc *xclient.XClient, ctx context.Context, typ, serviceMethod string, args *Args) {
	var reply int
	var err error
	switch typ {
	case "call":
		err = xc.Call(ctx, serviceMethod, args, &reply)
	case "broadcast":
		err = xc.Broadcast(ctx, serviceMethod, args, &reply)
	}
	if err != nil {
		log.Printf("%s %s error: %v", typ, serviceMethod, err)
	} else {
		log.Printf("%s %s success: %d + %d = %d", typ, serviceMethod, args.Num1, args.Num2, reply)
	}
}

/*
@Time    :   2022/11/11 16:54:40
@Author  :   victor2022
@Desc    :   发起调用
*/
func call(registry string) {
	d := xclient.NewGeeRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() { _ = xc.Close() }()
	// send request & receive response
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(xc, context.Background(), "call", "Foo.Sum", &Args{Num1: i, Num2: i * i})
		}(i)
	}
	wg.Wait()
}

func broadcast(registry string) {
	d := xclient.NewGeeRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func() { _ = xc.Close() }()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(xc, context.Background(), "broadcast", "Foo.Sum", &Args{Num1: i, Num2: i * i})
			// expect 2 - 5 timeout
			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
			foo(xc, ctx, "broadcast", "Foo.Sleep", &Args{Num1: i, Num2: i * i})
		}(i)
	}
	wg.Wait()
}

func main() {
	log.SetFlags(0)
	// 定义注册中心地址
	registryAddr := "http://localhost:9000/_geerpc_/registry0"
	// 开启注册中心
	var wg sync.WaitGroup
	wg.Add(1)
	go startRegistry(&wg)
	wg.Wait()

	time.Sleep(time.Second)

	// 开启服务
	wg.Add(2)
	go startServer(registryAddr, &wg)
	go startServer(registryAddr, &wg)
	wg.Wait()

	time.Sleep(time.Second)

	// 发起调用
	call(registryAddr)
	broadcast(registryAddr)

	ch := make(chan int)
	<-ch
}
