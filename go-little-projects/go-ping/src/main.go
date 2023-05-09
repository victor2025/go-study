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

	// 创建pingHandler
	pingHandler := ping.NewPingHandler(timeout, size, count, mode, address)

	// 开始ping操作
	pingHandler.StartPing()
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
