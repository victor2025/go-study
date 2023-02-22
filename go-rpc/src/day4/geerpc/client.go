/*
-*- encoding: utf-8 -*-
@File    :   client.go
@Time    :   2022/11/04 11:43:21
@Author  :   victor2022
@Version :   1.0
@Desc    :   client for geerpc
*/
package geerpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type Call struct {
	Seq           uint64
	ServiceMethod string      // 服务名和函数名
	Args          interface{} // 调用参数
	Reply         interface{} // 函数返回值
	Error         error
	Done          chan *Call // 在调用完成后用来通知客户端调用完成
}

/*
@Time    :   2022/11/04 11:47:55
@Author  :   victor2022
@Desc    :   通知调用方调用完成
*/
func (call *Call) done() {
	// 调用方通过channel接收调用完成的通知
	call.Done <- call
}

/*
@Time    :   2022/11/04 11:52:01
@Author  :   victor2022
@Desc    :   客户端主体
*/
type Client struct {
	cc       codec.Codec      // 编解码器
	opt      *Option          // 公有的option
	sending  sync.Mutex       // 保证消息有序发送
	header   codec.Header     // 消息头，同一个客户端可以复用
	mu       sync.Mutex       // 用于client自身操作时同步
	seq      uint64           // 请求编号，用于保证消息唯一性
	pending  map[uint64]*Call // 暂存未处理完成的请求，key为编号
	closing  bool             // 用户主动关闭
	shutdown bool             // 服务端要求停止，错误发生时被动关闭
}

var _ io.Closer = (*Client)(nil)

/*
@Time    :   2022/11/08 11:09:21
@Author  :   victor2022
@Desc    :   调用响应
*/
type clientResult struct {
	client *Client
	err    error
}

/*
@Time    :   2022/11/08 11:10:42
@Author  :   victor2022
@Desc    :   用来创建client的函数
*/
type newClientFunc func(conn net.Conn, opt *Option) (client *Client, err error)

// client已经被关闭时返回的默认错误
var ErrShutdown = errors.New("connection has been shutdown")

// Close implements io.Closer
func (client *Client) Close() error {
	client.mu.Lock()
	defer client.mu.Unlock()
	// 若已经关闭了
	if client.closing {
		return ErrShutdown
	}
	// 关闭
	client.closing = true
	return client.cc.Close()
}

/*
@Time    :   2022/11/04 11:55:58
@Author  :   victor2022
@Desc    :   返回当前客户端是否可用
*/
func (client *Client) IsAvailable() bool {
	client.mu.Lock()
	defer client.mu.Unlock()
	return !client.shutdown && !client.closing
}

/*
@Time    :   2022/11/04 12:02:41
@Author  :   victor2022
@Desc    :   注册调用到客户端中
*/
func (client *Client) registerCall(call *Call) (uint64, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	// 首先判断是否可用
	if client.closing || client.shutdown {
		return 0, ErrShutdown
	}
	// 设定调用的序列号，并将客户端序列号自增
	call.Seq = client.seq
	client.pending[call.Seq] = call
	client.seq++
	return call.Seq, nil
}

/*
@Time    :   2022/11/04 12:03:27
@Author  :   victor2022
@Desc    :   移除调用
*/
func (client *Client) removeCall(seq uint64) *Call {
	client.mu.Lock()
	defer client.mu.Unlock()
	call := client.pending[seq]
	delete(client.pending, seq)
	return call
}

/*
@Time    :   2022/11/04 12:05:12
@Author  :   victor2022
@Desc    :   客户端发生异常，主动终止调用
*/
func (client *Client) terminateCalls(err error) {
	client.sending.Lock()
	defer client.sending.Unlock()
	client.mu.Lock()
	defer client.mu.Unlock()
	client.shutdown = true
	// 所有尚未完成的调用都返回错误
	for _, call := range client.pending {
		call.Error = err
		call.done()
	}
}

/*
@Time    :   2022/11/04 12:15:22
@Author  :   victor2022
@Desc    :   接收响应
*/
func (client *Client) receive() {
	var err error
	// 循环接收响应，并判断是否出错
	for err == nil {
		var h codec.Header
		// 解析请求头
		if err = client.cc.ReadHeader(&h); err != nil {
			break
		}
		// 根据请求头中的序列号，移除暂存的响应
		call := client.removeCall(h.Seq)
		switch {
		case call == nil:
			// 该情况表明部分写失败，调用已经结束并移除了
			err = client.cc.ReadBody(nil)
		case h.Error != "":
			// 若存在错误，则手动结束
			call.Error = fmt.Errorf(h.Error)
			err = client.cc.ReadBody(nil)
			call.done()
		default:
			// 读取响应体
			err = client.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("read body error: " + err.Error())
			}
			// 完成响应
			call.done()
		}
	}
	// 跳出循环说明发生了错误，客户端发送终止信号
	client.terminateCalls(err)
}

/*
@Time    :   2022/11/04 12:27:51
@Author  :   victor2022
@Desc    :   创建新的客户端
*/
func NewClient(conn net.Conn, opt *Option) (*Client, error) {
	// 获取opt中定义的编解码器
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error: ", err)
		return nil, err
	}
	// 编码opt并以json的形式发送
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error: ", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

/*
@Time    :   2022/11/04 12:28:25
@Author  :   victor2022
@Desc    :   创建并启动客户端接收，内部方法
*/
func newClientCodec(cc codec.Codec, opt *Option) *Client {
	client := &Client{
		seq:     1, // 序列号从1开始
		cc:      cc,
		opt:     opt,
		pending: make(map[uint64]*Call),
	}
	// 开启client的循环接收
	go client.receive()
	return client
}

/*
@Time    :   2022/11/04 14:54:16
@Author  :   victor2022
@Desc    :   解析option
*/
func parseOptions(opts ...*Option) (*Option, error) {
	// 判断opt是否为空
	if len(opts) == 0 || opts[0] == nil {
		return DefaultOption, nil
	}
	// 若长度不为1，则超过了限制
	if len(opts) != 1 {
		return nil, errors.New("number of options is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = DefaultOption.MagicNumber
	// 根据配置选择合适的编解码器
	switch {
	case opt.CodecType == "" || opt.CodecType == "gob":
		opt.CodecType = DefaultOption.CodecType
	case opt.CodecType == "json":
		opt.CodecType = codec.JsonType
	}
	return opt, nil
}

/*
@Time    :   2022/11/04 14:58:00
@Author  :   victor2022
@Desc    :   创建和指定服务器的连接，带有超时处理
*/
func dialTimeout(f newClientFunc, network, address string, opts ...*Option) (client *Client, err error) {
	// 解析用户输入的配置
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}
	// 创建连接
	conn, err := net.DialTimeout(network, address, opt.ConnectionTimeout)
	if err != nil {
		return nil, err
	}
	// 若客户端创建失败的话，则关闭对应连接
	defer func() {
		if client == nil {
			_ = conn.Close()
		}
	}()

	// 异步创建client，通过channel返回结果
	ch := make(chan clientResult)
	go func() {
		client, err := f(conn, opt)
		ch <- clientResult{client: client, err: err}
	}()
	// 若不限制超时时间，则同步阻塞等待结果的返回
	if opt.ConnectionTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}
	// 若限制超时时间，则等待结果返回，若超时了也没有返回，则返回错误
	select {
	case <-time.After(opt.ConnectionTimeout):
		return nil, fmt.Errorf("rpc client: connect timeout: expect within %s", opt.ConnectionTimeout)
	case result := <-ch:
		return result.client, result.err
	}
}

/*
@Time    :   2022/11/08 11:26:15
@Author  :   victor2022
@Desc    :   包装了dialTimeout, 用来创建指定客户端
*/
func Dial(network, address string, opts ...*Option) (*Client, error) {
	return dialTimeout(NewClient, network, address, opts...)
}

/*
@Time    :   2022/11/04 15:05:39
@Author  :   victor2022
@Desc    :   发送请求
*/
func (client *Client) send(call *Call) {
	// 保证传输过程的线程安全
	client.sending.Lock()
	defer client.sending.Unlock()

	// 注册当前调用
	seq, err := client.registerCall(call)
	// 若发生错误，则该调用直接返回错误
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	// 准备请求头
	client.header.ServiceMethod = call.ServiceMethod
	client.header.Seq = seq
	client.header.Error = ""

	// 编码并发送请求
	if err := client.cc.Write(&client.header, call.Args); err != nil {
		// 若编码失败，则当前调用直接返回错误
		call := client.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

/*
@Time    :   2022/11/04 15:09:00
@Author  :   victor2022
@Desc    :   暴露给用户，异步执行请求，并返回Call结构
*/
func (client *Client) Go(serviceMethod string, args, reply interface{}, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}

	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done, // 调用结束标志
	}
	client.send(call)
	return call
}

/*
@Time    :   2022/11/04 15:10:53
@Author  :   victor2022
@Desc    :   暴露给用户，同步执行请求

	添加了context.Context作为传入参数，方便用户自定义进行超时处理
*/
func (client *Client) Call(ctx context.Context, serviceMethod string, args, reply interface{}) error {
	// 同步等待结果返回
	// 接收对应的call对象
	call := client.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	case <-ctx.Done():
		// ctx可以方便用户自定义调用何时超时
		// 通过context.WithTimeout方法定义超时时间
		client.removeCall(call.Seq)
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case call := <-call.Done:
		// Done为一个channel，在client.done()方法中传入当前call，表示调用完成
		return call.Error
	}
}
