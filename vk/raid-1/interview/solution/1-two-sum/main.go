package main

import (
	"fmt"
)

// Решение имеет линейную временную сложность O(n), а также линейную сложность по памяти
func twoSum(numbers []int, num int) bool {
	numbersMap := make(map[int]int)

	for i, el := range numbers {
		numbersMap[el] = i
	}

	for i, el := range numbers {
		findElement := num - el

		if n, ok := numbersMap[findElement]; ok && n != i {
			return true
		}
	}

	return false
}

func main() {
	fmt.Println(twoSum([]int{10, 15, 3, 7}, 17))

	fmt.Println(twoSum([]int{2, 2}, 4))

	fmt.Println(twoSum([]int{1, 2, 3, 4, 5, 6}, 12))
}
