package main

import "fmt"

func main() {
	var c1 = 'a'
	fmt.Printf("c1: %v--%c\n", c1, c1)
	var c2 = '1'
	fmt.Printf("c2: %v--%c\n", c2, c2)
	var c3 int = '啊'
	fmt.Printf("c3: %v--%c--%T\n", c3, c3, c3)
}
