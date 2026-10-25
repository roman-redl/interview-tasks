package main

import (
	"fmt"
)

func main() {
	s := "asdq"
	// s[0] = 'r' // compilation error
	println(s)

	slices()
	printString()
}

func slices() {
	a := []string{"a", "b", "c"} // [a, b, c] len:3 cap:3
	b := a[1:2]                  // [b] len:1 cap:2
	b[0] = "q"                   // [q] len:1 cap:2
	b = append(b, "d")           // [q, d] len:2 cap:2
	b = append(b, "e")           // [q, d, e] len:3 cap:4
	b[0] = "x"
	fmt.Println(a) // [a, q, d]
}

func slices2() {
	a := []string{"a", "b", "c"} // [a, b, c] len:3 cap:3
	b := a[2:2]                  // [] len:0 cap:2
	b[0] = "q"                   // panic out of range
	b = append(b, "d")
	b = append(b, "e")
	b[0] = "x"
	fmt.Println(a)
}

func slices3() {
	a := []string{"a", "b", "c"} // [a, b, c] len:3 cap:3
	b := a[1:4]                  // panic: runtime error: slice bounds out of range [:4] with capacity 3
	b[0] = "q"
	b = append(b, "d")
	b = append(b, "e")
	b[0] = "x"
	fmt.Println(a)
}

func printString() {
	s := "Юadsf"

	for i := range s {
		fmt.Println(i) // 0 2 3 4 5
	}
}
