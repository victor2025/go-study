package main

import (
	"flag"
	"go-trans/handlers"
)

var (
	help   bool
	sMode  bool
	addr   string
	input  string
	output string
	port   string
)

func main() {
	getCmdArgs()
	if help {
		flag.Usage()
		return
	}
	// 判断参数合法性
	if sMode {
		sHandler := handlers.NewServerHandler(port, output)
		sHandler.Handle()
	} else {
		cHandler := handlers.NewClientHandler(addr, port, input)
		cHandler.Handle()
	}
}

// 获取命令行参数
func getCmdArgs() {
	// 使用flag包读取命令行参数
	flag.BoolVar(&help, "h", false, "打印本帮助文档")
	flag.BoolVar(&sMode, "r", false, "接收模式")
	flag.StringVar(&output, "o", ".received/", "接收模式，文件保存路径")
	flag.StringVar(&addr, "s", "localhost", "发送模式，目标地址")
	flag.StringVar(&input, "i", "", "发送模式，要发送的文件")
	flag.StringVar(&port, "p", "20235", "接受/发送模式，端口号")
	flag.Parse()
}
