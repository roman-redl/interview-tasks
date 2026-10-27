package main

import "fmt"

// Sum solve 2
func Sum(pir [][]int, n int) int {
	sum := 0
	for i := 0; i <= n; i++ {
		sum += pir[n][i]
	}

	return sum
}

func main() {
	pyramid := [][]int{{1}, {3, 5}, {7, 9, 11}}
	fmt.Println(Sum(pyramid, 2))
}
