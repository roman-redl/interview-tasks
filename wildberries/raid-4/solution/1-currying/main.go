package main

import (
	"fmt"
)

func add(a int) func(b int) int {
	return func(b int) int {
		return a + b
	}
}

func main() {
	d := add(3)(5)
	fmt.Println(d) // 8
}
