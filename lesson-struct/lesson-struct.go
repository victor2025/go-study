package main

import (
	"fmt"
	"math"
)

type Position struct {
	x, y float64
}

// 专属于Position的方法
func (p Position) Dis0() float64 {
	return CalcDis(p.x, p.y)
}

func Dis1(p Position) float64 {
	return CalcDis(p.x, p.y)
}

func CalcDis(x, y float64) float64 {
	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

func main() {
	// 创建结构体
	pos := Position{100, 50}
	// 使用.调用方法
	fmt.Println(pos.Dis0())
	// 使用普通函数
	fmt.Println(Dis1(pos))
}
