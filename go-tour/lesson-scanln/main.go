package main

import (
	"fmt"
	"os"
)

func main() {
	var name string
	var age int
	// 读取数据
	fmt.Print("Input your name please: ")
	fmt.Scanln(&name)
	fmt.Print("Input your age please: ")
	fmt.Scanln(&age)
	fmt.Printf("Hello! %s, your age is %d\n", name, age)
	fmt.Println(os.Environ())
}
