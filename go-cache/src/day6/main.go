/*
-*- encoding: utf-8 -*-
@File    :   main.go
@Time    :   2022/10/26 16:30:34
@Author  :   victor2022
@Version :   1.0
@Desc    :   Main for geecache
*/
package main

import (
	"flag"
	"fmt"
	"geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

/*
@Time    :   2022/10/28 20:27:15
@Author  :   victor2022
@Desc    :   创建group
*/
func createGroup() *geecache.Group {
	return geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

/*
@Time    :   2022/10/28 20:27:46
@Author  :   victor2022
@Desc    :   开始提供服务，注册分布式节点
*/
func startCacheServer(addr string, addrs []string, gee *geecache.Group) {
	// 创建服务器
	peers := geecache.NewHTTPPool(addr)
	// 设置分布式节点
	peers.Set(addrs...)
	// 注册分布式节点
	gee.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

/*
@Time    :   2022/10/28 20:31:49
@Author  :   victor2022
@Desc    :   提供服务调用api
*/
func startAPIServer(apiAddr string, gee *geecache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := gee.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func main() {
	var port int
	var api bool
	// 获取启动参数
	flag.IntVar(&port, "port", 8001, "Geecache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	// 创建提供查询的score group
	gee := createGroup()
	if api {
		// 通过协程启动一个暴露api的服务
		go startAPIServer(apiAddr, gee)
	}
	// 开启分布式缓存，每个节点中都注册了其他节点的信息
	startCacheServer(addrMap[port], []string(addrs), gee)
}
