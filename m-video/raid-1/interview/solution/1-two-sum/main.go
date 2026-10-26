package main

import "fmt"

// twoSum за один проход: храним в map уже виденные значения (value ->
// index); для очередного числа ищем в map пару target - v.
// Сложность: O(n) времени, O(n) памяти; brute force двумя циклами был бы
// O(n^2) времени / O(1) памяти.
func twoSum(nums []int, target int) []int {
	seen := make(map[int]int, len(nums)) // value -> index
	for i, v := range nums {
		if j, ok := seen[target-v]; ok {
			return []int{j, i}
		}
		seen[v] = i
	}
	return nil
}

func main() {
	fmt.Println(twoSum([]int{2, 7, 11, 15}, 9)) // [0 1]
	fmt.Println(twoSum([]int{3, 2, 4}, 6))      // [1 2]
	fmt.Println(twoSum([]int{3, 3}, 6))         // [0 1] — дубликат обрабатывается: второй 3 находит первый в map
}
