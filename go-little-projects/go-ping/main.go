package main

import (
	"flag"
	"os"
	"ping"
)

var (
	help    bool
	timeout int64
	size    int
	count   int
	mode    bool
	address string
)

func main() {
	// 获取参数
	getCmdArgs()
	// 打印帮助文档
	if help {
		flag.Usage()
		return
	}

	pingHandler := ping.NewPingHandler(timeout, size, count, mode, address)

	pingHandler.StartPing()
	// // 开始ping操作
	// // 和目标地址建立连接
	// conn, err := net.DialTimeout("ip:icmp", address, time.Duration(timeout)*time.Millisecond)
	// u.CheckErr(err, func() { os.Exit(1) })

	// // 创建连接成功后，开始处理
	// // 处理完成后关闭连接
	// defer conn.Close()

	// // 获取字节形式的icmp报文
	// msg := icmp.GetBytes(size)

	// // 通过连接发送报文
	// conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	// _, err = conn.Write(*msg)
	// u.CheckErr(err, func() { os.Exit(1) })

	// // 读响应
	// buf := make([]byte, 65535)
	// n, err := conn.Read(buf)
	// u.CheckErr(err, func() { os.Exit(1) })
	// ipMsg, err := p.ParseBytes2Ip(buf[0:n])
	// u.CheckErr(err, func() {})

	// // 处理响应
	// log.Printf("Host %v received response from %v: size=%dBytes ttl=%dms\n",
	// 	ipMsg.Aim, ipMsg.Source, ipMsg.Size-28, ipMsg.TTL)

}

// 获取命令行参数
func getCmdArgs() {
	// 使用flag包读取命令行参数
	flag.BoolVar(&help, "h", false, "打印本帮助文档")
	flag.Int64Var(&timeout, "w", 1000, "请求超时时长，单位毫秒")
	flag.IntVar(&size, "l", 32, "请求发送缓冲区大小")
	flag.IntVar(&count, "n", 4, "发送请求的数目")
	flag.BoolVar(&mode, "t", false, "是否持续发送(测试模式)")
	flag.Parse()
	// 解析目标地址
	address = os.Args[len(os.Args)-1]
}
