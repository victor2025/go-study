/*
-*- encoding: utf-8 -*-
@File    :   lru.go
@Time    :   2022/10/25 11:28:01
@Author  :   victor2022
@Version :   1.0
@Desc    :   LRU strategy for cache
*/
package lru

import "container/list"

type (
	/*
		@Time    :   2022/10/25 11:30:47
		@Author  :   victor2022
		@Desc    :   cache 的主体结构
	*/
	Cache struct {
		// 允许使用的最大内存
		maxBytes int64
		// 当前已经使用的内存
		nBytes int64
		// 双向链表结构
		ll *list.List
		// 值是双向链表中对应节点的指针
		cache map[string]*list.Element
		// 当一个元素完全移除时的回调函数
		OnEvicted func(key string, value Value)
	}

	/*
		@Time    :   2022/10/25 11:31:57
		@Author  :   victor2022
		@Desc    :   元素结构
	*/
	entry struct {
		key   string
		value Value
	}

	/*
		@Time    :   2022/10/25 11:32:47
		@Author  :   victor2022
		@Desc    :   值对象
	*/
	Value interface {
		// Len()方法返回存储值大小
		Len() int
	}
)

/*
@Time    :   2022/10/25 11:39:01
@Author  :   victor2022
@Desc    :   创建新的Cache对象
*/
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

/*
@Time    :   2022/10/25 11:41:20
@Author  :   victor2022
@Desc    :   查找元素，找到之后将元素移动到队尾
*/
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 由key找到对应的
	if ele, ok := c.cache[key]; ok {
		// 移动到队尾
		c.ll.MoveToBack(ele)
		// 找到对应的Value
		// 将当前类型转换位entry类型
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

/*
@Time    :   2022/10/25 14:39:29
@Author  :   victor2022
@Desc    :   根据淘汰策略移除元素
*/
func (c *Cache) RemoveOldest() {
	// 弹出队首元素
	ele := c.ll.Front()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		// 在map中删除对应元素
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		// 若存在回调函数，则调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

/*
@Time    :   2022/10/25 14:44:27
@Author  :   victor2022
@Desc    :   添加元素操作
*/
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 若当前元素已经存在，则将元素移动到队尾
		c.ll.MoveToBack(ele)
		// 更新当前元素
		kv := ele.Value.(*entry)
		// 更新大小
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 若当前元素不存在
		// 创建并放入当前元素，push之后，会返回一个指向该元素的指针
		ele := c.ll.PushBack(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 淘汰元素，0表示无上限
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

/*
@Time    :   2022/10/25 15:01:47
@Author  :   victor2022
@Desc    :   获取缓存中数据的数目
*/
func (c *Cache) Len() int {
	return c.ll.Len()
}
