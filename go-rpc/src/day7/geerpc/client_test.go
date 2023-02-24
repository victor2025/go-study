/*
-*- encoding: utf-8 -*-
@File    :   client_test.go
@Time    :   2022/11/09 15:41:52
@Author  :   victor2022
@Version :   1.0
@Desc    :   test for client of geerpc
*/
package geerpc

import (
	"context"
	"net"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

/*-----------测试连接创建超时--------------*/

/*
@Time    :   2022/11/09 15:52:37
@Author  :   victor2022
@Desc    :   测试连接超时处理
*/
func TestClient_dialTimeout(t *testing.T) {
	t.Parallel()
	l, _ := net.Listen("tcp", ":0")
	// 一个需要耗时两秒的函数
	f := func(conn net.Conn, opt *Option) (client *Client, err error) {
		_ = conn.Close()
		time.Sleep(time.Second * 2)
		return nil, nil
	}
	// 同时运行两个测试
	// 超时时间为1s
	t.Run("timeout", func(t *testing.T) {
		_, err := dialTimeout(f, "tcp", l.Addr().String(), &Option{ConnectionTimeout: time.Second})
		_assert(err != nil && strings.Contains(err.Error(), "connect timeout"), "expect a timeout error")
	})
	// 不规定超时时间
	t.Run("0", func(t *testing.T) {
		_, err := dialTimeout(f, "tcp", l.Addr().String(), &Option{ConnectionTimeout: 0})
		_assert(err == nil, "0 means no limit")
	})
}

/*-----------测试服务调用超时--------------*/

type Bar int

// 被调用的函数，需要2s的运行时间
func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(time.Second * 2)
	return nil
}

func startServer(addr chan string) {
	var b Bar
	_ = Register(&b)
	// pick a free port
	l, _ := net.Listen("tcp", ":0")
	addr <- l.Addr().String()
	Accept(l)
}

/*
@Time    :   2022/11/09 15:53:24
@Author  :   victor2022
@Desc    :   测试调用
*/
func TestClient_Call(t *testing.T) {
	t.Parallel()
	addrCh := make(chan string)
	go startServer(addrCh)
	addr := <-addrCh
	time.Sleep(time.Second)
	// 调用，包含1s远程调用超时
	t.Run("client timeout", func(t *testing.T) {
		client, _ := Dial("tcp", addr)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		var reply int
		err := client.Call(ctx, "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), ctx.Err().Error()), "expect a timeout error")
	})
	// 调用，包含1s服务端调用超时
	t.Run("server handle timeout", func(t *testing.T) {
		client, _ := Dial("tcp", addr, &Option{
			HandleTimeout: time.Second,
		})
		var reply int
		err := client.Call(context.Background(), "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), "handle timeout"), "expect a timeout error")
	})
}

/*
@Time    :   2022/11/11 16:32:36
@Author  :   victor2022
@Desc    :   测试基于不同协议客户端的创建
*/
func TestXDial(t *testing.T) {
	if runtime.GOOS == "linux" {
		ch := make(chan struct{})
		// 本地socket
		addr := "/tmp/geerpc.sock"
		// 采用默认服务器监听socket
		go func() {
			_ = os.Remove(addr)
			l, err := net.Listen("unix", addr)
			if err != nil {
				t.Fatal("failed to listen unix socket")
			}
			ch <- struct{}{}
			Accept(l)
		}()
		<-ch
		// 创建基于unix socket协议的rpc client
		_, err := XDial("unix@" + addr)
		_assert(err == nil, "failed to connect unix socket")
	}
}
