package main

import (
	"fmt"
)

// 比较方法的值传递和引用传递
type Position struct {
	x, y float64
	name string
}

func (p Position) Change0() {
	p.x = -p.x
	p.y = -p.y
	p.name = "changed by 0"
}

func (p *Position) Change1() {
	p.x = -p.x
	p.y = -p.y
	p.name = "changed by 1"
}

func (p Position) print() {
	fmt.Printf("%f %f %s", p.x, p.y, p.name)
}

func main() {
	pos := Position{
		100,
		100,
		"original name",
	}
	// 测试值传递
	pos.Change0()
	pos.print()
	fmt.Println()
	// 测试引用传递
	pos.Change1()
	pos.print()
}
