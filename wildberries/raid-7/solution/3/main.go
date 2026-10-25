package main

import "fmt"

func main() {
	s := "Юadsf"

	for i := range s {
		fmt.Println(i) // 0 2 3 4 5
	}
}
