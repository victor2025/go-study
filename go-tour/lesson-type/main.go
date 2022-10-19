package main

import (
	"fmt"
	"strconv"
	"unsafe"
)

func Check(i interface{}) {
	// 直接打印类型
	fmt.Printf("type of i is: %T\n", i)
	fmt.Printf("size of i is: %d\n", unsafe.Sizeof(i))
	fmt.Printf("default value of i is: %v\n", i)
	// 通过.type获取类型
	switch i.(type) {
	case int:
		fmt.Println("i is a int")
	case string:
		fmt.Println("i is a string")
	default:
		fmt.Println("i is another type")
	}
}

func main() {
	Check(123)
	Check("hello")
	Check(nil)
	// conv
	n1 := strconv.FormatInt(123, 10)
	fmt.Printf("n1: %v:%T\n", n1, n1)
	n2 := strconv.FormatFloat(1.2, 'f', 10, 64)
	fmt.Printf("n2: %v:%T\n", n2, n2)
}
