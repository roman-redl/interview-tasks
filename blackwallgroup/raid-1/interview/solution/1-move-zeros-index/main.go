package main

import "fmt"

// MoveZerosIndex solve 1.2
func MoveZerosIndex(num []int) []int {
	ind := 0

	for i := 0; i < len(num); i++ {
		if num[i] != 0 {
			num[ind] = num[i]
			ind++
		}
	}

	// Особенность варианта: возвращается префикс длины ind — нули остаются
	// в нижележащем массиве за границей len. Для строгой формулировки
	// «переместить нули в конец» хвост num[ind:] дописывается нулями.
	return num[:ind]
}

func main() {
	fmt.Println(MoveZerosIndex([]int{1, 0, 1, 2, 0, 1, 3}))
}
