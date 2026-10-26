package main

import (
	"fmt"
)

//func a() {
//	x := []int{}
//	x = append(x, 0)
//	fmt.Println(cap(x), len(x))
//
//	x = append(x, 1)
//	fmt.Println(cap(x), len(x))
//
//	x = append(x, 2)
//	fmt.Println(cap(x), len(x))
//
//	y := append(x, 3, 7, 8)
//
//	fmt.Println(cap(x), len(x))
//	fmt.Println(x, y)
//
//	z := append(x, 4)
//	z = append(x, 5)
//	fmt.Println(cap(x), len(x))
//
//	fmt.Println(x, y, z)
//}

// [0 1 2] [0 1 2 4] [0 1 2 4]
func a() {
	x := []int{}
	x = append(x, 0)
	x = append(x, 1)
	x = append(x, 2)
	y := append(x, 3)
	z := append(x, 4)
	fmt.Println(x, y, z)
}

func main() {
	a()
}
