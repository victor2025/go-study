package protocols

import (
	"errors"
	"fmt"
	"utils"
)

type IP struct {
	Size   uint16
	TTL    uint8
	Source IpAddr
	Aim    IpAddr
	Body   []byte
	SeqNum uint16
}

func ParseBytes2Ip(data []byte) (*IP, error) {
	// 校验data长度
	l := len(data)
	if l < 20 {
		return nil, errors.New("data length invalid")
	}
	// 解析size
	size := uint16(l)
	ttl := uint8(utils.Bytes2Uint64(data[8:9]))
	source := data[12:16]
	aim := data[16:20]
	body := data[20:]
	seqNum := uint16(utils.Bytes2Uint64(body[6:8])) + 1
	return &IP{
		Size:   size,
		TTL:    ttl,
		Source: source,
		Aim:    aim,
		Body:   body,
		SeqNum: seqNum,
	}, nil
}

type IpAddr []byte

func (addr IpAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", addr[0], addr[1], addr[2], addr[3])
}
