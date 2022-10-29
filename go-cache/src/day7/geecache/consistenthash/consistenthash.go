/*
-*- encoding: utf-8 -*-
@File    :   consistenthash.go
@Time    :   2022/10/27 22:36:39
@Author  :   victor2022
@Version :   1.0
@Desc    :   consistent hash for geecache
*/
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/*
@Time    :   2022/10/27 22:37:19
@Author  :   victor2022
@Desc    :   定义hash函数形式
*/
type Hash func(data []byte) uint32

/*
@Time    :   2022/10/27 22:39:37
@Author  :   victor2022
@Desc    :   一致性哈希主结构
*/
type Map struct {
	hash     Hash           // 哈希算法，可以自定义
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环
	hashMap  map[int]string // 虚拟节点与真实节点的映射关系
}

/*
@Time    :   2022/10/27 22:42:09
@Author  :   victor2022
@Desc    :   创建新的Map结构
*/
func New(replicas int, fn Hash) *Map {
	// 创建新的map结构
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	// 若hash函数为空，则使用默认hash函数
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

/*
@Time    :   2022/10/27 22:47:57
@Author  :   victor2022
@Desc    :   添加真实节点，并生成对应的虚拟节点
*/
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 生成虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 排序所有的虚拟节点
	sort.Ints(m.keys)
}

/*
@Time    :   2022/10/27 22:52:20
@Author  :   victor2022
@Desc    :   获取当前key对应的最近的节点
*/
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 生成当前key的hash值
	hash := int(m.hash([]byte(key)))
	// 二分查找大于当前hash值的最近虚拟节点
	// 当找不到时，返回数组长度
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 获取真实节点名
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
