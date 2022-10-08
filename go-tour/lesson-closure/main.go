package main

import "fmt"

// 闭包练习

func main() {
	// 普通闭包
	multi2 := compute(2)
	fmt.Println(multi2(2))
	multi3 := compute(3)
	fmt.Println(multi3(2))
	// 斐波那契数列
	// 第一个闭包
	f1 := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Print(f1(), " ")
	}
	fmt.Println()
	// 第二个闭包
	f2 := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Print(f2(), " ")
	}
	fmt.Println()
}

/*
*
闭包生成函数
*/
func compute(c int) func(a int) int {
	return func(a int) int {
		res := c * a
		return res
	}
}

func fibonacci() func() int {
	c1 := 0
	c2 := 1
	return func() int {
		res := c1 + c2
		c1 = c2
		c2 = res
		return res
	}
}
