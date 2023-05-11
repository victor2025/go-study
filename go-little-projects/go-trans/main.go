package main

import (
	"flag"
	"log"
	"go-trans/handlers"
)

var (
	help      bool
	isReceive bool = false
	addr      string
	input     string
	output    string
	port      string
)

func main() {
	getCmdArgs()
	if help {
		flag.Usage()
		return
	}

	log.Println("--- go-trans: a file transmitter by go ---")
	// 判断参数合法性
	if isReceive {
		sHandler := handlers.NewReceiveHandler(port, output)
		sHandler.Handle()
	} else {
		cHandler := handlers.NewSendHandler(addr, port, input)
		cHandler.Handle()
	}
}

// 获取命令行参数
func getCmdArgs() {
	// 使用flag包读取命令行参数
	flag.BoolVar(&help, "h", false, "print this help doc 打印本帮助文档")
	flag.StringVar(&output, "o", ".received/", "file output dir 文件保存路径")
	flag.StringVar(&addr, "s", "", "send mode (receive mode for default) [aim address] 发送模式(默认为接收模式) [目标地址]")
	flag.StringVar(&input, "i", "", "file for send 要发送的文件")
	flag.StringVar(&port, "p", "20235", "port number 端口号")
	flag.Parse()
	// 根据目标地址是否为空，修改模式
	if len(addr) == 0 {
		isReceive = true
	}
}
