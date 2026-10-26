package main

import (
	"fmt"
	"sort"
)

// По убыванию значений ключи уже не помощники — сортировать надо пары:
// map разворачивается в слайс структур (ключ едет рядом со своим
// значением), сортируется по значению в обратном порядке. Для устойчивого
// порядка при равных значениях — вторичный компаратор по ключу.
type pair struct {
	k, v int
}

func printDescValues(m map[int]int) {
	pairs := make([]pair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].v != pairs[j].v {
			return pairs[i].v > pairs[j].v
		}
		return pairs[i].k > pairs[j].k
	})

	for _, p := range pairs {
		fmt.Printf("%d %d\n", p.k, p.v)
	}
}

func main() {
	printDescValues(map[int]int{1: 10, 2: 5, 3: 30, 10: 1})
}
