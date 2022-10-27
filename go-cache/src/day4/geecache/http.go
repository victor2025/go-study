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
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_geecache/"

/*
@Time    :   2022/10/26 16:10:41
@Author  :   victor2022
@Desc    :   使用该结构存储对应的访问路径
*/
type HttpPool struct {
	self     string
	basePath string
}

/*
@Time    :   2022/10/26 16:16:29
@Author  :   victor2022
@Desc    :   创建对应指定请求的连接池
*/
func NewHTTPPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: defaultBasePath,
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
