/*
-*- encoding: utf-8 -*-
@File    :   getter_test.go
@Time    :   2022/10/25 16:08:14
@Author  :   victor2022
@Version :   1.0
@Desc    :   test for geecache.Getter
*/
package geecache

import (
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}
