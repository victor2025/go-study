/*
-*- encoding: utf-8 -*-
@File    :   http.go
@Time    :   2022/10/26 16:06:41
@Author  :   victor2022
@Version :   1.0
@Desc    :   HttpServer for geecache
*/

package geecache

import (
	"fmt"
	"geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_geecache/"
	defaultReplicas = 50
)

/*
@Time    :   2022/10/26 16:10:41
@Author  :   victor2022
@Desc    :   使用该结构存储对应的访问路径

	实现了PeerPicker接口，可以从节点列表中获取对应的远程节点调用器
*/
type HttpPool struct {
	self        string                 // 自身路径
	basePath    string                 // 基本路径
	mu          sync.Mutex             // 负责在获取匹配项或者获取httpGetter时进行同步
	peers       *consistenthash.Map    // 一致性哈希，根据
	httpGetters map[string]*httpGetter // 远程节点与httpGetter的映射
}

/*
@Time    :   2022/10/28 12:45:35
@Author  :   victor2022
@Desc    :   远程数据获取器
*/
type httpGetter struct {
	baseURL string
}

/*
@Time    :   2022/10/26 16:16:29
@Author  :   victor2022
@Desc    :   创建对应指定请求的连接池
*/
func NewHTTPPool(self string) *HttpPool {
	return &HttpPool{
		self:        self,
		basePath:    defaultBasePath,
		peers:       nil,
		httpGetters: nil,
	}
}

/*
@Time    :   2022/10/26 16:18:26
@Author  :   victor2022
@Desc    :   日志打印
*/
func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

/*
@Time    :   2022/10/26 16:20:13
@Author  :   victor2022
@Desc    :   HTTP服务端必须实现的方法
*/
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 校验路径，是否是从指定域名访问的
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// 解析参数
	// /<basePath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	// 校验参数数目
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// 分配参数
	groupName := parts[0]
	key := parts[1]

	// 从缓存中获取对应的group
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	// 从group中获取key
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置响应头并返回信息
	// 响应体标志为字节流
	w.Header().Set("Content-Type", "application/octet-stream")
	// 写入字节信息
	w.Write(view.ByteSlice())
}

/*
@Time    :   2022/10/28 12:48:20
@Author  :   victor2022
@Desc    :   从远程节点获取对应数据
*/
func (h *httpGetter) Get(group string, key string) ([]byte, error) {
	// 拼接远程节点的调用数据
	// QueryEscape()方法实现网址编码
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	// 远程调用其他节点
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	// 若请求成功，则在方法返回前关闭响应体
	defer res.Body.Close()
	// 校验响应状态
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returns: %v", res.Status)
	}
	// 读取数据
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}
	// 返回数据
	return bytes, nil
}

/*
@Time    :   2022/10/28 13:09:38
@Author  :   victor2022
@Desc    :   向pool中添加分布式节点，并进行初始化
*/
func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

/*
@Time    :   2022/10/28 13:13:38
@Author  :   victor2022
@Desc    :   实现PeerPicker方法，按照key获取对应的远程节点调用者
*/
func (p *HttpPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// 获取对应的peer，并进行调用
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		// 对该节点进行调用
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}
