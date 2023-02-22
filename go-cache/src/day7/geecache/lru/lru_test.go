/*
-*- encoding: utf-8 -*-
@File    :   lru_test.go
@Time    :   2022/10/25 15:03:33
@Author  :   victor2022
@Version :   1.0
@Desc    :   test for lru
*/
package lru

import (
	"reflect"
	"testing"
)

type String string

/*
@Time    :   2022/10/25 15:04:25
@Author  :   victor2022
@Desc    :   重写Len方法
*/
func (d String) Len() int {
	return len(d)
}

/*
@Time    :   2022/10/25 15:07:36
@Author  :   victor2022
@Desc    :   测试取值
*/
func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

/*
@Time    :   2022/10/25 15:07:50
@Author  :   victor2022
@Desc    :   测试值的淘汰
*/
func TestRemoveoldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Add(k1, String(v1))
	lru.Add(k2, String(v2))
	lru.Add(k3, String(v3))

	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

/*
@Time    :   2022/10/25 15:09:34
@Author  :   victor2022
@Desc    :   测试回调函数的调用
*/
func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("123456"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}
