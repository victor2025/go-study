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

func main() {
	geecache.NewGroup("scores", 2<<10, geecache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	// curl http://localhost:9000/_geecache/scores/Tom
	// curl http://localhost:9000/_geecache/scores/kkk
	addr := "localhost:9000"
	// 建立面向指定地址请求的服务
	peers := geecache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
