/*
-*- encoding: utf-8 -*-
@File    :   gob.go
@Time    :   2022/11/03 17:07:17
@Author  :   victor2022
@Version :   1.0
@Desc    :   encode and decode by gob of golang object for rpc communication
*/
package codec

import (
	"bufio"        // 实现了带有缓冲的io
	"encoding/gob" // go官方库中定义的编解码方法
	"io"
	"log"
)

/*
@Time    :   2022/11/03 17:14:13
@Author  :   victor2022
@Desc    :   Gob编解码主结构
*/
type GobCodec struct {
	conn io.ReadWriteCloser
	buf  *bufio.Writer
	dec  *gob.Decoder // gob解码器，持有io流，可以实现从流中读取并解码
	enc  *gob.Encoder // gob编码器，持有io流，可以实现将数据编码后写入流
}

/*
@Time    :   2022/11/03 17:16:46
@Author  :   victor2022
@Desc    :   gob编解码器构造函数，属于NewCodecFunc类型
*/
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	// 创建新的gob编解码器
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

/*
@Time    :   2022/11/03 17:18:16
@Author  :   victor2022
@Desc    :   Codec接口定义的读Header的方法
*/
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

/*
@Time    :   2022/11/03 17:21:21
@Author  :   victor2022
@Desc    :   Codec接口定义的读Body的方法
*/
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	// 一次性写入buf中所有的数据，并主动关闭io流
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	// 尝试编码
	if err := c.enc.Encode(h); err != nil {
		log.Println("rpc codec: gob error encoding header: ", err)
		return err
	}
	if err := c.enc.Encode(body); err != nil {
		log.Println("rpc codec: gob error encoding body: ", err)
		return err
	}
	return nil
}

/*
@Time    :   2022/11/03 17:30:58
@Author  :   victor2022
@Desc    :   定义Close方法
*/
func (c *GobCodec) Close() error {
	return c.conn.Close()
}

var _ Codec = (*GobCodec)(nil)
