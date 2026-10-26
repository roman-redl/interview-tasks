package main

import (
	"fmt"
	"sort"
)

// mergeMeetings объединяет пересекающиеся (и стыкующиеся) интервалы:
// сортировка по началу, затем один линейный проход со склейкой.
// O(n log n) на сортировку + O(n) на проход; результат — новый слайс,
// вход не мутируется.
func mergeMeetings(meetings [][2]int) [][2]int {
	if len(meetings) == 0 {
		return nil
	}

	sorted := append([][2]int(nil), meetings...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i][0] < sorted[j][0] })

	merged := make([][2]int, 0, len(sorted))
	cur := sorted[0]
	for _, m := range sorted[1:] {
		if m[0] <= cur[1] { // пересекается или стыкуется вплотную
			if m[1] > cur[1] {
				cur[1] = m[1]
			}
			continue
		}
		merged = append(merged, cur)
		cur = m
	}
	return append(merged, cur)
}

func main() {
	fmt.Println(mergeMeetings([][2]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}})) // [[1 6] [8 10] [15 18]]
	fmt.Println(mergeMeetings([][2]int{{1, 4}, {4, 5}}))                     // [[1 5]] — стык [1,4]+[4,5]
	fmt.Println(mergeMeetings([][2]int{{5, 6}, {1, 2}, {3, 4}}))             // [[1 2] [3 4] [5 6]] — вход не отсортирован
	fmt.Println(mergeMeetings([][2]int{{1, 10}, {2, 3}, {4, 5}}))            // [[1 10]] — вложенные поглощаются
	fmt.Println(mergeMeetings(nil))                                          // []
}
