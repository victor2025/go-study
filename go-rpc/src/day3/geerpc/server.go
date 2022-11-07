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
	"errors"
	"geerpc/codec"
	"go/ast"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
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
type Server struct {
	serviceMap sync.Map // 存放所有注册了的服务
}

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
@Time    :   2022/11/07 15:20:04
@Author  :   victor2022
@Desc    :   向当前服务器中注册服务
*/
func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	// 向server中放入当前service
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc server: service already defined" + s.name)
	}
	return nil
}

/*
@Time    :   2022/11/07 15:23:57
@Author  :   victor2022
@Desc    :   向默认服务器中注册服务
*/
func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

/*
@Time    :   2022/11/07 15:25:57
@Author  :   victor2022
@Desc    :   根据服务和方法名找到对应的服务
*/
func (server *Server) findService(serviceMethod string) (svc *service, mType *methodType, err error) {
	// 找到分隔符的位置
	dotIdx := strings.LastIndex(serviceMethod, ".")
	if dotIdx < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	// 分割出服务名和方法名
	serviceName, methodName := serviceMethod[:dotIdx], serviceMethod[dotIdx+1:]
	// 找出对应的服务对象
	svci, ok := server.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	// 找出对应的方法
	svc = svci.(*service)
	mType = svc.method[methodName]
	if mType == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}

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
	// 向函数中传入对应的编码解码器
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
		// 解码请求
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
		// 使用goroutine进行异步处理
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
@Desc    :   rpc请求结构体，包含当前请求的所有信息
*/
type request struct {
	h            *codec.Header // 请求头
	argv, replyv reflect.Value // 请求中的参数和对应的响应
	mType        *methodType   // 当前请求调用的方法类型
	svc          *service      // 当前请求对应的服务对象
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
	// 找到当前请求对应的服务
	req.svc, req.mType, err = server.findService(h.ServiceMethod)
	if err != nil {
		return req, err
	}
	// 创建请求参数
	req.argv = req.mType.newArgv()
	req.replyv = req.mType.newReplyv()

	// 保证argvi是一个指针，用来接收反序列化后的请求体
	argvi := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	// 读取输入参数的值
	if err = cc.ReadBody(argvi); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}
	return req, nil
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
@Time    :   2022/11/03 22:28:21
@Author  :   victor2022
@Desc    :   处理请求
*/
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	// 完成之后解锁一次
	defer wg.Done()
	// 调用service进行调用
	err := req.svc.call(req.mType, req.argv, req.replyv)
	// 若出错，则返回出错信息
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, req.h, invalidRequest, sending)
		return
	}
	// 若没有出错，则返回调用结果
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
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
@Time    :   2022/11/07 11:12:55
@Author  :   victor2022
@Desc    :   被调用的方法类型，各个部分都由反射类型表示
*/
type methodType struct {
	method    reflect.Method // 被调用的方法
	ArgType   reflect.Type   // 第一个参数，可能是普通变量或者指针
	ReplyType reflect.Type   // 第二个参数，一定为指针
	numCalls  uint64         // 调用的次数
}

/*
@Time    :   2022/11/07 11:15:13
@Author  :   victor2022
@Desc    :   获取调用数目
*/
func (m *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&m.numCalls)
}

/*
@Time    :   2022/11/07 11:16:18
@Author  :   victor2022
@Desc    :   创建新的参数
*/
func (m *methodType) newArgv() reflect.Value {
	var argv reflect.Value
	// 根据变量是否为指针参数进行不同的创建
	if m.ArgType.Kind() == reflect.Ptr {
		// 创建指针类型
		argv = reflect.New(m.ArgType.Elem())
	} else {
		// 创建普通类型
		argv = reflect.New(m.ArgType).Elem()
	}
	return argv
}

/*
@Time    :   2022/11/07 11:23:26
@Author  :   victor2022
@Desc    :   创建新的replyv变量
*/
func (m *methodType) newReplyv() reflect.Value {
	// reply必须为指针类型
	replyv := reflect.New(m.ReplyType.Elem())
	// 若为map或slice，则需要进行初始化
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		// 初始化反射map
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		// 初始化反射切片
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return replyv
}

/*
@Time    :   2022/11/07 11:57:31
@Author  :   victor2022
@Desc    :   服务单元，一个service对应一个提供服务的结构体
*/
type service struct {
	name   string                 // 映射的结构体名称
	typ    reflect.Type           // 结构体类型
	rcvr   reflect.Value          // 结构体实例在反射对象
	method map[string]*methodType // 映射结构体所有的方法
}

/*
@Time    :   2022/11/07 12:14:03
@Author  :   victor2022
@Desc    :   将结构体注册为新的服务
*/
func newService(rcvr interface{}) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	s.typ = reflect.TypeOf(rcvr)
	// 查看方法是否暴露
	if !ast.IsExported(s.name) {
		log.Fatalf("rpc server: %s is not a valid service name", s.name)
	}
	// 注册当前结构体中的方法
	s.registerMethods()
	return s
}

/*
@Time    :   2022/11/07 14:55:03
@Author  :   victor2022
@Desc    :   注册service对应的结构体中的方法
*/
func (s *service) registerMethods() {
	s.method = make(map[string]*methodType)
	// 遍历获取结构体中的方法，并注册
	for i := 0; i < s.typ.NumMethod(); i++ {
		method := s.typ.Method(i)
		mType := method.Type
		// 跳过输入输出参数数目不满足条件的方法，0参数为自身，其他参数分别为arg和reply
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		// 跳过输出参数类型不满足条件的方法
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		// 跳过输出参数不满足条件的方法
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		// 所有条件都满足，创建methodType注册对应的方法
		s.method[method.Name] = &methodType{
			method:    method,
			ArgType:   argType,
			ReplyType: replyType,
		}
		log.Printf("rpc server: register %s.%s\n", s.name, method.Name)
	}
}

/*
@Time    :   2022/11/07 14:52:37
@Author  :   victor2022
@Desc    :   检测参数是否为暴露的类型或者为当前可以访问到的类型
*/
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

/*
@Time    :   2022/11/07 15:01:22
@Author  :   victor2022
@Desc    :   通过service调用方法
*/
func (s *service) call(m *methodType, argv, replyv reflect.Value) error {
	// 将当前方法被调用的次数+1
	atomic.AddUint64(&m.numCalls, 1)
	// 获取当前方法的反射对象
	f := m.method.Func
	// 通过方法的反射对象对指定结构体的该方法进行调用，其中所用的所有参数均属于反射空间
	returnValues := f.Call([]reflect.Value{s.rcvr, argv, replyv})
	// 判断调用是否出错
	if errInter := returnValues[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	return nil
}
