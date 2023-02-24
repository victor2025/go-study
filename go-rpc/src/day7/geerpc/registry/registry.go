/*
-*- encoding: utf-8 -*-
@File    :   registry.go
@Time    :   2023/02/23 12:39:41
@Author  :   victor2022
@Version :   1.0
@Desc    :   main body of registry
*/
package registry

import (
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// 注册中心结构体
// 用来注册和获取服务
// 可以从列表中删除失效服务
type GeeRegistry struct {
	timeout time.Duration          // 超时时间
	mu      sync.Mutex             // 保护服务列表
	servers map[string]*ServerItem // 存储服务实例
}

// 服务实例
type ServerItem struct {
	Addr  string    // 服务对应地址
	start time.Time // 开始服务时间
}

const (
	defaultPath    = "/_geerpc_/registry" // 默认访问路径
	defaultTimeout = time.Minute * 5      // 默认超时时间
)

/*
@Time    :   2023/02/23 16:15:03
@Author  :   victor2022
@Desc    :   创建一个注册中心实例

	传入服务列表的过期时间
*/
func New(timeout time.Duration) *GeeRegistry {
	return &GeeRegistry{
		servers: make(map[string]*ServerItem),
		timeout: timeout,
	}
}

var DefaultGeeRegistry = New(defaultTimeout) // 默认注册中心

/*
@Time    :   2023/02/23 16:18:12
@Author  :   victor2022
@Desc    :   向指定注册中心注册服务
*/
func (r *GeeRegistry) putServer(addr string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// 获取对应的服务
	s := r.servers[addr]
	if s == nil {
		// 若对应的服务实例不存在，则创建
		r.servers[addr] = &ServerItem{
			Addr:  addr,
			start: time.Now(),
		}
	} else {
		// 若当前服务已存在，则更新服务的开始时间
		s.start = time.Now()
	}
}

/*
@Time    :   2023/02/23 16:18:59
@Author  :   victor2022
@Desc    :   获取指定注册中心中存活的服务列表
*/
func (r *GeeRegistry) aliveServer() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	var alive []string
	for addr, s := range r.servers {
		if r.timeout == 0 || s.start.Add(r.timeout).After(time.Now()) {
			// 若没有设定超时时间或者没有超过超时时间，则表明当前服务存活
			alive = append(alive, addr)
		} else {
			// 否则删除该服务
			delete(r.servers, addr)
		}
	}
	// 按照字典序排序服务列表
	sort.Strings(alive)
	return alive
}

// 采用http协议接收请求
func (r *GeeRegistry) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("processing request through", req.Host)
	switch req.Method {
	// 服务的信息放在请求头中
	case "GET":
		// GET用来获取服务地址
		// 在响应头中放所有存活服务的地址
		w.Header().Set("X-Geerpc-Servers", strings.Join(r.aliveServer(), ","))
	case "POST":
		// POST请求用来注册服务
		// 从请求头中获取地址
		addr := req.Header.Get("X-Geerpc-Server")
		if addr == "" {
			// 若地址为空，则返回错误
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		r.putServer(addr)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// 为注册中心添加一个处理路径
func (r *GeeRegistry) HandleHTTP(registryPath string) {
	// 由于GeeRegistry属于HttpHandler，因此可以直接被用作Handler
	http.Handle(registryPath, r)
	log.Println("rpc registry path:", registryPath)
}

// 默认的服务器开启方法
func HandleHTTP() {
	DefaultGeeRegistry.HandleHTTP(defaultPath)
}

// 心跳函数，服务向注册中心发送心跳
func Heartbeat(registry, addr string, duration time.Duration) {
	if duration == 0 {
		// 保证有足够的时间发送心跳包
		duration = defaultTimeout - time.Duration(1)*time.Minute
	}
	var err error
	err = sendHeartbeat(registry, addr)
	go func() {
		// 开始计时，创建计时器
		t := time.NewTicker(duration)
		for err == nil {
			// 等待计时结束
			<-t.C
			// 发送下一次心跳
			err = sendHeartbeat(registry, addr)
		}
	}()
}

// 服务向注册中心发送心跳
func sendHeartbeat(registry, addr string) error {
	log.Println(addr, " send heart beat to registry ", registry)
	httpClient := &http.Client{}
	// 创建请求，向注册中心发送心跳
	req, _ := http.NewRequest("POST", registry, nil)
	// 设置请求头，在其中包含自己的身份
	req.Header.Set("X-Geerpc-Server", addr)
	// 开始请求
	if _, err := httpClient.Do(req); err != nil {
		log.Println("rpc server: heart beat err:", err)
		return err
	}
	return nil
}
