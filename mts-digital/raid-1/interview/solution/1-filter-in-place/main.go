package main

import "fmt"

// [1,2,3,4,5] -> [1,3,5]
// [2,2,2] -> []
// [1,2,2] -> [1]
func filter(a []int) []int {
	var n int

	for i := 0; i < len(a); i++ {
		if a[i]%2 != 0 {
			a[n] = a[i]
			n++
		}
	}

	return a[:n]
}

func main() {
	fmt.Println(filter([]int{2, 2, 1, 4, 4, 5, 7, 7, 7}))
}
