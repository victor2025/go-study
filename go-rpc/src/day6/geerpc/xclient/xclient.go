/*
-*- encoding: utf-8 -*-
@File    :   xclient.go
@Time    :   2022/12/03 14:39:08
@Author  :   victor2022
@Version :   1.0
@Desc    :   client who supports load balance
*/
package xclient

import (
	"context"
	. "geerpc"
	"io"
	"reflect"
	"sync"
)

// 支持负载均衡的客户端
type XClient struct {
	d       Discovery
	mode    SelectMode
	opt     *Option
	mu      sync.Mutex
	clients map[string]*Client
}

var _ io.Closer = (*XClient)(nil)

func NewXClient(d Discovery, mode SelectMode, opt *Option) *XClient {
	return &XClient{
		d:       d,
		mode:    mode,
		opt:     opt,
		clients: make(map[string]*Client),
	}
}

/*
@Time    :   2022/12/03 14:46:53
@Author  :   victor2022
@Desc    :   关闭xclient
*/
func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, client := range xc.clients {
		// 暂时不处理错误
		_ = client.Close()
		delete(xc.clients, key)
	}
	return nil
}

/*
@Time    :   2022/12/03 14:56:22
@Author  :   victor2022
@Desc    :   与传入的地址创建连接
*/
func (xc *XClient) dial(rpcAddr string) (*Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	// 从客户端缓存中取出对应的client
	client, ok := xc.clients[rpcAddr]
	// 检验当前地址对应的client是否可用
	if ok && !client.IsAvailable() {
		_ = client.Close()
		delete(xc.clients, rpcAddr)
		client = nil
	}
	// 若当前client为空，则创建对应地址的的client
	if client == nil {
		var err error
		// 创建对应地址的client
		client, err = XDial(rpcAddr, xc.opt)
		if err != nil {
			return nil, err
		}
		// 将client存入本地缓存中
		xc.clients[rpcAddr] = client
	}
	return client, nil
}

/*
@Time    :   2022/12/03 14:58:42
@Author  :   victor2022
@Desc    :   向指定服务发起调用
*/
func (xc *XClient) call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// 按照调用地址查找对应的客户端
	client, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	// 发起调用
	return client.Call(ctx, serviceMethod, args, reply)
}

/*
@Time    :   2022/12/03 15:01:38
@Author  :   victor2022
@Desc    :   提供给用户的带有负载均衡算法的远程调用函数
*/
func (xc *XClient) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// 首先尝试通过负载均衡算法从发现器中获取对应的服务地址
	rpcAddr, err := xc.d.Get(xc.mode)
	if err != nil {
		return err
	}
	// 向获取到的服务地址发送调用请求
	return xc.call(rpcAddr, ctx, serviceMethod, args, reply)
}

/*
@Time    :   2022/12/03 15:03:43
@Author  :   victor2022
@Desc    :   向服务列表中所有的服务广播调用命令
*/
func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex // 保证错误和函数调用的线程安全
	var e error
	replyDone := reply == nil // 用来判断是否需要返回值，若不需要，则无需等待返回
	// 将context包装为带有cancel方法的context
	ctx, cancel := context.WithCancel(ctx)
	for _, rpcAddr := range servers {
		// 加锁一次
		wg.Add(1)
		go func(rpcAddr string) {
			// 完成后解锁一次
			defer wg.Done()
			// 使用reply的副本接收调用结果
			var clonedReply interface{}
			if reply != nil {
				clonedReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			// 发起调用
			err := xc.call(rpcAddr, ctx, serviceMethod, args, clonedReply)
			// 处理返回
			mu.Lock()
			if err != nil && e == nil {
				e = err
				cancel() // 若某个调用出现问题，则取消所有的未完成的调用
			}
			if err == nil && !replyDone {
				// 将cloneReply的值赋给reply
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(clonedReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	return e
}
