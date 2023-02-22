/*
-*- encoding: utf-8 -*-
@File    :   server.go
@Time    :   2022/11/03 17:33:15
@Author  :   victor2022
@Version :   1.0
@Desc    :   Server of geerpc
*/
package geerpc

import (
	"encoding/json"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

// 魔数，作为框架标记
const MagicNumber = 0x3bef5c

/*
@Time    :   2022/11/03 17:35:20
@Author  :   victor2022
@Desc    :   数据传输中的消息的配置项
*/
type Option struct {
	MagicNumber int
	CodecType   codec.Type
}

/*
@Time    :   2022/11/03 17:35:04
@Author  :   victor2022
@Desc    :   定义默认Option
*/
var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// 服务端结构体
type Server struct{}

/*
@Time    :   2022/11/03 17:40:08
@Author  :   victor2022
@Desc    :   服务端构造函数
*/
func NewServer() *Server {
	return &Server{}
}

// 默认Server，单例
var DefaultServer = NewServer()

/*
@Time    :   2022/11/03 17:42:12
@Author  :   victor2022
@Desc    :   循环等待接收并创建连接

	创建连接后，将连接交给另一个go-routine进行处理
*/
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		go server.ServeConn(conn)
	}
}

/*
@Time    :   2022/11/03 17:44:39
@Author  :   victor2022
@Desc    :   采用默认服务器接收请求
*/
func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

/*
@Time    :   2022/11/04 10:35:56
@Author  :   victor2022 
@Desc    :   循环处理连接
*/
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	// 最后一定要关闭io
	defer func() {
		_ = conn.Close()
	}()
	// 开始处理
	var opt Option
	// 将json解析为Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error:", err)
		return
	}
	// 比对MagicNumber
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	// 解析其他部分
	// 获取编解码器
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}
	// 对剩余信息进行处理
	server.serveCodec(f(conn))
}

// 当返回值发生问题时，使用该变量作为作为占位符
var invalidRequest = struct{}{}

/*
@Time    :   2022/11/03 17:58:11
@Author  :   victor2022
@Desc    :   解码请求并发送给函数进行请求
*/
func (server *Server) serveCodec(cc codec.Codec) {
	sending := new(sync.Mutex) // 同步器，保证传输过程不被干扰
	wg := new(sync.WaitGroup)  // 同步器，保证所有的请求都被处理
	// 持续读取请求
	for {
		req, err := server.readRequest(cc)
		if err != nil {
			if req == nil {
				break // 解析失败时终止循环，否则针对其他错误进行响应
			}
			// 返回针对错误的响应
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		// 加锁，等待请求响应
		wg.Add(1)
		// 处理请求，使用sending锁保证单个连接中的多个报文逐个发送
		go server.handleRequest(cc, req, sending, wg)
	}
	// 等待所有响应都处理完成
	wg.Wait()
	// 关闭io流
	_ = cc.Close()
}

/*
@Time    :   2022/11/03 22:14:09
@Author  :   victor2022
@Desc    :   rpc请求结构体
*/
type request struct {
	h            *codec.Header // 请求头
	argv, replyv reflect.Value // 请求中的参数和对应的响应
}

/*
@Time    :   2022/11/03 22:19:29
@Author  :   victor2022
@Desc    :   读请求头
*/
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

/*
@Time    :   2022/11/03 22:20:41
@Author  :   victor2022
@Desc    :   读取请求
*/
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	// TODO 假定当前请求体为string类型
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

/*
@Time    :   2022/11/03 22:25:10
@Author  :   victor2022
@Desc    :   发送响应
*/
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	// 加锁，避免相互干扰
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

/*
@Time    :   2022/11/03 22:28:21
@Author  :   victor2022
@Desc    :   处理请求
*/
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// 完成之后解锁一次
	defer wg.Done()
	// TODO简单处理
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}
