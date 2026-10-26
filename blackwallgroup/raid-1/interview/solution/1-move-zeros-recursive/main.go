package main

import "fmt"

// MoveZerosRecursive solve 1.1
func MoveZerosRecursive(zeros []int) []int {
	return MoveZero(zeros, 0)
}

func MoveZero(zeros []int, n int) []int {
	if n == len(zeros)-1 {
		return zeros
	} else if zeros[n] == 0 {
		if zeros[n+1] == 0 {
			zeros = MoveZero(zeros, n+1)
		}
		zeros[n], zeros[n+1] = zeros[n+1], zeros[n]
	}

	return MoveZero(zeros, n+1)
}

func main() {
	fmt.Println(MoveZerosRecursive([]int{1, 0, 1, 2, 0, 1, 3}))
}
