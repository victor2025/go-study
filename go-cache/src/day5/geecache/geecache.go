/*
-*- encoding: utf-8 -*-
@File    :   geecache.go
@Time    :   2022/10/25 15:35:39
@Author  :   victor2022
@Version :   1.0
@Desc    :   main body of geecache

	使用group包装了cache
	可为group指定对应的值获取方法，若本地缓存中没有目标值，则根据值获取方法获取目标值
*/
package geecache

import (
	"fmt"
	"log"
	"sync"
)

/*
@Time    :   2022/10/25 15:36:23
@Author  :   victor2022
@Desc    :   用于用户实现的Getter，在对应数据找不到的时候调用其中的Get方法

	该接口属于函数式接口，其中只有一个方法。
	在使用时既可以以函数的方式出现，又可以以(继承了当前方法的)结构体的方式出现
	与Java的函数式接口类比，Java中也通过函数式接口避免了创建新的对象
*/
type Getter interface {
	Get(key string) ([]byte, error)
}

/*
@Time    :   2022/10/25 15:37:27
@Author  :   victor2022
@Desc    :   获取数据的函数结构，用于实现Getter
*/
type GetterFunc func(key string) ([]byte, error)

/*
@Time    :   2022/10/25 16:17:24
@Author  :   victor2022
@Desc    :   一个group对应一个命名空间和对应的数据获取方法
*/
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker // 每一个group对应一个Picker负责远程调用
}

var (
	// 同步量
	mu sync.RWMutex
	// 存放group的map结构
	groups = make(map[string]*Group)
)

/*
@Time    :   2022/10/25 15:38:42
@Author  :   victor2022
@Desc    :   实现了Getter
*/
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/*
@Time    :   2022/10/25 16:21:10
@Author  :   victor2022
@Desc    :   创建group的函数
*/
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	// 同步
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	// 存入map中
	groups[name] = g
	return g
}

/*
@Time    :   2022/10/25 16:24:33
@Author  :   victor2022
@Desc    :   由group名称获取对应的group
*/
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

/*
@Time    :   2022/10/25 16:34:10
@Author  :   victor2022
@Desc    :   从group中获取值的操作
*/
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 从主缓存中获取数据
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[GeeCache] hit key: %q", key)
		return v, nil
	}
	// 命中失败则调用load方法
	return g.load(key)
}

/*
@Time    :   2022/10/25 16:36:53
@Author  :   victor2022
@Desc    :   调用其他方法获取数据
*/
func (g *Group) load(key string) (ByteView, error) {
	// 若存在其他节点，则先尝试从其他节点获取当前值
	if g.peers != nil {
		// 从group中获取映射节点名称
		if peer, ok := g.peers.PickPeer(key); ok {
			// 尝试从对应节点获取结果
			if value, err := g.getFromPeer(peer, key); err == nil {
				return value, nil
			} else {
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}
	}
	// 从本地方法获取
	return g.getLocally(key)
}

/*
@Time    :   2022/10/27 22:17:22
@Author  :   victor2022
@Desc    :   由本地获取对应值
*/
func (g *Group) getLocally(key string) (ByteView, error) {
	// 获取数据
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	// 只读
	value := ByteView{b: cloneBytes(bytes)}
	// 将获得的数据添加到缓存中
	g.populateCache(key, value)
	return value, nil
}

/*
@Time    :   2022/10/27 22:18:17
@Author  :   victor2022
@Desc    :   将值放入本地缓存中
*/
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

/*
@Time    :   2022/10/28 13:26:16
@Author  :   victor2022
@Desc    :   注册peers到当前group
*/
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

/*
@Time    :   2022/10/28 13:32:55
@Author  :   victor2022
@Desc    :   从远程节点获取值
*/
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	bytes, err := peer.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}
