package main

import "fmt"

func main() {
	// 创建map
	tempMap := make(map[string]string)
	// 放入数据
	tempMap["李恒威"] = "victor2022"
	fmt.Println(tempMap["李恒威"])
	// 多值获取
	val, exist := tempMap["111"]
	if exist {
		fmt.Println(val)
	} else {
		fmt.Println("val is Nil")
	}
}
