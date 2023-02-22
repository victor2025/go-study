/*
-*- encoding: utf-8 -*-
@File    :   codec.go
@Time    :   2022/11/03 11:11:42
@Author  :   victor2022
@Version :   1.0
@Desc    :   encoder and decoder of geerpc
*/
package codec

import "io"

type Header struct {
	ServiceMethod string // 请求的服务名和方法，Service.Method
	Seq           uint64 // 请求序列号
	Error         string
}

/*
@Time    :   2022/11/03 11:28:16
@Author  :   victor2022
@Desc    :   抽象编解码接口，具体编解码方式可以自定义实现
*/
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

/*
@Time    :   2022/11/03 11:29:24
@Author  :   victor2022
@Desc    :   定义编解码器的构造函数格式
*/
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

/*
@Time    :   2022/11/03 11:33:28
@Author  :   victor2022
@Desc    :   提供的序列化类型
*/
const (
	GobType  Type = "application/gob"  // gob类型
	JsonType Type = "application/json" // Json类型
)

// 编解码器map，用于保存和查找编解码器
var NewCodecFuncMap map[Type]NewCodecFunc

/*
@Time    :   2022/11/05 23:07:35
@Author  :   victor2022
@Desc    :   初始化函数，编译器能够保证init函数在main函数之前执行
*/
func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	// Gob编解码器创建函数
	NewCodecFuncMap[GobType] = NewGobCodec
	// Json编解码器创建函数
	NewCodecFuncMap[JsonType] = NewJsonCodec
}
