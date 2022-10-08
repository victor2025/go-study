package main

import (
	"fmt"
	"math"
)

// 定义接口，空接口可以对应任意类型(类似Java的Object类)
type ICalc interface {
	Dis() float64
}

type Position struct {
	x, y float64
}

// 隐式实现
func (p Position) Dis() float64 {
	return math.Sqrt(p.x*p.x + p.y*p.y)
}

func main() {
	// 接口测试
	var p ICalc
	p1 := Position{100, 100}
	p2 := Position{100, 100}

	p = p1
	fmt.Println(p.Dis())
	Show(p)

	p = &p2
	fmt.Println(p.Dis())
	Show(p)

	// 类型断言
	val, ok := p.(ICalc)
	fmt.Println(val, ok)
}

func Show(i ICalc) {
	fmt.Printf("The result of i is: %f\n", i.Dis())
}
