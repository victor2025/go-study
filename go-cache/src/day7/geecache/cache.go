/*
-*- encoding: utf-8 -*-
@File    :   cache.go
@Time    :   2022/10/25 15:25:01
@Author  :   victor2022
@Version :   1.0
@Desc    :   None
*/
package geecache

import (
	"geecache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

/*
@Time    :   2022/10/25 15:27:17
@Author  :   victor2022
@Desc    :   线程安全的缓存添加算法
*/
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 若没有初始化lru，则新建一个lru
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

/*
@Time    :   2022/10/25 15:29:19
@Author  :   victor2022
@Desc    :   线程安全的缓存获取算法
*/
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 若lru还没有初始化，则直接返回空值
	if c.lru == nil {
		return
	}
	// 若能够找到对应的数据，则返回
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
