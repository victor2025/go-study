package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"go-trans/utils"
	"net"
)

const (
	magic              = 151 // magic number
	headSize           = 4   // size of head
	ByteType TransType = 0
	StrType  TransType = 1
	NumType  TransType = 2
	EndType  TransType = 3
)

type TRANS struct {
	Head    *THead
	Content []byte
}

type THead struct {
	Magic     uint8
	Type      TransType
	TotalSize uint16
}

type TransType uint8

func ByteTransMsg(data []byte) *TRANS {
	// 配置数据
	size := uint16(headSize + len(data)) // head + content
	head := &THead{
		Magic:     magic,
		Type:      ByteType,
		TotalSize: size,
	}
	return &TRANS{
		Head:    head,
		Content: data,
	}
}

func StrTransMsg(data string) *TRANS {
	trans := ByteTransMsg([]byte(data))
	trans.Head.Type = StrType
	return trans
}

func NumTransMsg(data int64) *TRANS {
	bytes := make([]byte, 64)
	n := binary.PutVarint(bytes, data)
	trans := ByteTransMsg(bytes[:n])
	trans.Head.Type = NumType
	return trans
}

func EndTransMsg() *TRANS {
	trans := ByteTransMsg([]byte{})
	trans.Head.Type = EndType
	return trans
}

func ReceiveNextTrans(conn net.Conn) (*TRANS, error) {
	var err error
	// read head
	headBytes := make([]byte, headSize)
	_, err = conn.Read(headBytes)
	utils.HandleError(err)
	if err != nil {
		return nil, err
	}
	head, err := parseHead(headBytes)
	utils.HandleError(err)

	// read filename
	content := make([]byte, head.TotalSize-headSize)
	_, err = conn.Read(content)
	utils.HandleError(err, func() {})
	return &TRANS{
		Head:    head,
		Content: content,
	}, nil
}

// 解析头部
func parseHead(bytes []byte) (*THead, error) {
	if len(bytes) != headSize {
		return nil, errors.New("invalid head byte data")
	}
	m := uint8(bytes[0])
	if m != magic {
		return nil, fmt.Errorf("invalid magic number, expect %d, got %d", magic, m)
	}
	t := uint8(bytes[1])
	size := utils.Bytes2Uint64(bytes[2:4])
	return &THead{
		Magic:     uint8(magic),
		Type:      TransType(t),
		TotalSize: uint16(size),
	}, nil
}

func (t *TRANS) Bytes() []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, t.Head)
	buf.Write(t.Content)
	return buf.Bytes()
}
