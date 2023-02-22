/*
-*- encoding: utf-8 -*-
@File    :   byteview.go
@Time    :   2022/10/25 15:19:31
@Author  :   victor2022
@Version :   1.0
@Desc    :   readonly struct for cache
*/
package geecache

type ByteView struct {
	b []byte
}

/*
@Time    :   2022/10/25 15:20:43
@Author  :   victor2022
@Desc    :   获取当前结构体的大小
*/
func (v ByteView) Len() int {
	return len(v.b)
}

/*
@Time    :   2022/10/25 15:22:45
@Author  :   victor2022
@Desc    :   将当前数据转化为一个新的byte切片
*/
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

/*
@Time    :   2022/10/25 15:23:23
@Author  :   victor2022
@Desc    :   将当前对象存储数据转化为string输出
*/
func (v ByteView) String() string {
	return string(v.b)
}

/*
@Time    :   2022/10/25 15:24:28
@Author  :   victor2022
@Desc    :   克隆一个byte数组，保证数据的只读
*/
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
