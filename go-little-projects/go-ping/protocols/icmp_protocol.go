package protocols

import (
	"bytes"
	"encoding/binary"
)

// icmp报文结构体
type ICMP struct {
	Type     uint8
	Code     uint8
	CheckSum uint16
	ID       uint16
	SeqNum   uint16
}

func GetICMPingMsg() *ICMP {
	return &ICMP{
		Type:     8,
		Code:     0,
		CheckSum: 0,
		ID:       0,
		SeqNum:   0,
	}
}

func (icmp *ICMP) GetBytes(size int) *[]byte {
	// 写入byte数组
	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	buffer.Write(make([]byte, size))
	bytes := buffer.Bytes()
	// 校验和
	checkSum := checkSum(bytes)
	bytes[2] = byte(checkSum >> 8)
	bytes[3] = byte(checkSum)
	return &bytes
}

// icmp的校验和算法
func checkSum(data []byte) uint16 {
	length := len(data)
	idx := 0
	var sum uint32
	for idx < length {
		sum += uint32(data[idx])<<8 + uint32(data[idx+1])
		idx += 2
		if (idx + 1) >= length {
			break
		}
	}
	if idx+1 == length {
		sum += uint32(data[idx])
	}
	// 取高16位
	hi16 := sum >> 16
	for hi16 != 0 {
		sum = hi16 + uint32(uint16(sum)) // 高低16位相加
		hi16 = sum >> 16
	}
	return uint16(^sum)
}
