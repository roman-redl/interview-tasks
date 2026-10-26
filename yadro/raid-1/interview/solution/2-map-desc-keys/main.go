package main

import (
	"fmt"
	"sort"
)

// Map в Go неупорядочен (порядок обхода специально рандомизирован):
// чтобы вывести пары по убыванию ключей, ключи собираются в слайс,
// сортируются по убыванию, и по ним печатается map.
func printDescKeys(m map[int]int) {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	for _, k := range keys {
		fmt.Printf("%d %d\n", k, m[k])
	}
}

func main() {
	printDescKeys(map[int]int{1: 1, 2: 2, 3: 3, 10: 10})
}
