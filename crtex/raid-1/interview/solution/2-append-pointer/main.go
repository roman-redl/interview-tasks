package main

import (
	"fmt"
	"sort"
)

func main() {
	v := []int{3, 4, 1, 2, 5}

	ap(&v)
	sr(v)

	fmt.Println(v)
}

func ap(arr *[]int) {
	*arr = append(*arr, 10)
}

func sr(arr []int) {
	sort.Ints(arr)
}
