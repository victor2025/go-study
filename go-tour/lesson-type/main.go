package main

import "fmt"

func Check(i interface{}) {
	switch i.(type) {
	case int:
		fmt.Println("i is a int")
	case string:
		fmt.Println("i is a string")
	default:
		fmt.Println("i is a other type")
	}
}

func main() {
	Check(123)
	Check("hello")
	Check(nil)
}
