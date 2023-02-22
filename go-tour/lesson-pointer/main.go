/*
* -*- encoding: utf-8 -*-
@File    :   main.go
@Time    :   2022/10/19 15:21:50
@Author  :   victor2022
@Version :   1.0
@Desc    :   None
*
*/
package main

import (
	"fmt"
	"unsafe"
)

func main() {
	n1 := 1
	fmt.Printf("n1: %v\n", n1)
	p1 := &n1
	fmt.Printf("p1: %v, %T\n", p1, p1)
	*p1++
	fmt.Printf("n1: %v\n", n1)
	fmt.Printf("p1: %v, %T\n", p1, p1)
	a := false
	fmt.Printf("a: %v, %T, %d \n", a, a, unsafe.Sizeof(a))
}
