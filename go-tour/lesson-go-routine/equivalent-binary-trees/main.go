package main

import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk 步进 tree t 将所有的值从 tree 发送到 channel ch。
func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	// 中序遍历
	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same 检测树 t1 和 t2 是否含有相同的值。
func Same(t1, t2 *tree.Tree) bool {
	ch1, ch2 := make(chan int), make(chan int)
	// 执行walk
	go func() {
		Walk(t1, ch1)
		close(ch1)
	}()
	go func() {
		Walk(t2, ch2)
		close(ch2)
	}()
	// 开始遍历
	for {
		n1, n2 := <-ch1, <-ch2
		if n1 == 0 && n2 == 0 {
			return true
		} else if n1 != n2 {
			return false
		}
	}
}

func main() {
	// 建树
	t1 := tree.New(1)
	// ch在关闭后只能读取到0值
	ch := make(chan int)
	// 测试Walk
	go func() {
		Walk(t1, ch)
		close(ch)
	}()
	for {
		num := <-ch
		if num == 0 {
			break
		}
		fmt.Print(num, " ")
	}
	fmt.Println()
	// 测试Same
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}
