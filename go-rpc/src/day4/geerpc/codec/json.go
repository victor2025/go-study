/*
-*- encoding: utf-8 -*-
@File    :   json.go
@Time    :   2022/11/05 17:53:40
@Author  :   victor2022
@Version :   1.0
@Desc    :   encode and decode by json
*/
package codec

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

/*
@Time    :   2022/11/05 17:57:50
@Author  :   victor2022
@Desc    :   json编解码器主结构
*/
type JsonCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *json.Decoder // json解码器
	enc  *json.Encoder // jso编码器
}

/*
@Time    :   2022/11/05 22:54:45
@Author  :   victor2022
@Desc    :   创建新的Json编解码器
*/
func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	// 创建新的json编解码器
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

// Close implements Codec
func (c *JsonCodec) Close() error {
	return c.conn.Close()
}

// ReadBody implements Codec
func (c *JsonCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// ReadHeader implements Codec
func (c *JsonCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// Write implements Codec
func (c *JsonCodec) Write(h *Header, body interface{}) (err error) {
	// 若写入出错，则主动关闭io流
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	// 编码并写入bufio
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: json error encoding header: ", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: json error encoding body: ", err)
		return err
	}
	return
}

var _ Codec = (*JsonCodec)(nil)
