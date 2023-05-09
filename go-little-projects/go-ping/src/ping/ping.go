package ping

import (
	"fmt"
	"math"
	"net"
	"os"
	p "protocols"
	"time"
	u "utils"
)

type PingHandler struct {
	icmp       *p.ICMP // icmp ping请求byte数组
	start      time.Time
	timeout    time.Duration
	size       int
	count      int
	addr       string
	sendCnt    int
	successCnt int
	totalTime  float64
	minTime    float32
	maxTime    float32
}

func NewPingHandler(timeout int64, size, count int, mode bool, addr string) *PingHandler {
	// 解析mode
	if mode {
		count = math.MaxInt
	}
	// 返回对象
	return &PingHandler{
		icmp:    p.GetICMPingMsg(),
		timeout: time.Duration(timeout) * time.Millisecond,
		size:    size,
		count:   count,
		addr:    addr,
		minTime: math.MaxFloat32,
		maxTime: 0.0,
	}
}

func (h *PingHandler) StartPing() {
	// 创建连接
	conn, err := net.DialTimeout("ip:icmp", h.addr, h.timeout)
	u.CheckErr(err, func() {
		os.Exit(1)
	})
	defer conn.Close()
	fmt.Printf("Ping %v(%v) from %v:\n", h.addr, conn.RemoteAddr(), conn.LocalAddr())
	// 开始ping操作
	h.start = time.Now()
	for i := 0; i < h.count; i++ {
		start := time.Now()

		// 设置本次超时时间
		conn.SetDeadline(start.Add(h.timeout))

		// 发送报文
		_, err := conn.Write(*h.icmp.GetBytes(h.size))
		u.CheckErr(err)
		h.icmp.IncrSeqNum() // 增加序列号
		h.sendCnt++

		// 接收响应
		buf := make([]byte, 65535)
		n, err := conn.Read(buf)
		u.CheckErr(err)
		ipMsg, err := p.ParseBytes2Ip(buf[:n])
		u.CheckErr(err, nil, func() {
			// 获取持续时间
			h.successPing(ipMsg, time.Since(start))
		})
		// 等待一秒
		time.Sleep(time.Duration(1) * time.Second)
	}
	// 完成后处理
	h.endPing()
}

func (h *PingHandler) successPing(res *p.IP, dur time.Duration) {
	if res == nil {
		fmt.Printf("received invalid response from\n")
		return
	}
	// 生成时间
	durTime := float32(dur.Microseconds()) / 1000
	// 记录状态
	h.successCnt++
	h.totalTime += float64(durTime)
	if durTime > h.maxTime {
		h.maxTime = durTime
	}
	if durTime < h.minTime {
		h.minTime = durTime
	}
	// 打印本次结果
	fmt.Printf("received %d bytes from %v: icmp_seq=%d ttl=%d rtt=%.2fms\n",
		res.Size-28, res.Source, res.SeqNum, res.TTL, durTime)
}

func (h *PingHandler) endPing() {
	lostPercent := float32(h.sendCnt-h.successCnt) / float32(h.sendCnt) * 100
	avgTime := float32(h.totalTime / float64(h.successCnt))
	fmt.Printf("\n--- Ping %v statistics ---\n", h.addr)
	fmt.Printf("sent %d packs, received %d packs, lost %.2f%%, cost %dms\n",
		h.sendCnt, h.successCnt, lostPercent, time.Since(h.start).Milliseconds())
	fmt.Printf("avg rtt: %.2fms, min rtt: %.2fms, max rtt: %.2fms\n",
		avgTime, h.minTime, h.maxTime)
}
