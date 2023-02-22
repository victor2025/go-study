/*
-*- encoding: utf-8 -*-
@File    :   codec.go
@Time    :   2022/11/03 11:11:42
@Author  :   victor2022
@Version :   1.0
@Desc    :   encoder and decoder of geeprc
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
@Desc    :   编解码器的构造函数
*/
type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

/*
@Time    :   2022/11/03 11:33:28
@Author  :   victor2022
@Desc    :   相关常量
*/
const (
	GobType  Type = "application/gob" // 自定义go对象类型
	JsonType Type = "application/json"
)

// 编解码器map，用于保存和查找编解码器
var NewCodecFuncMap map[Type]NewCodecFunc

func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec // TODO
}
