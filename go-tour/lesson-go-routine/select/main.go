package main

import (
	"fmt"
	"time"
)

func main() {
	tick := time.Tick(100 * time.Millisecond)
	ring := time.After(1000 * time.Millisecond)
	// 采用select选择不同channel进行接收
	for {
		select {
		case <-tick:
			fmt.Print("tick ")
		case <-ring:
			fmt.Print("\nBOOM!!!")
			return
		default:
			fmt.Print(". ")
			time.Sleep(50 * time.Millisecond)
		}
	}
}
