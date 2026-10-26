package main

import "fmt"

// longestSubstring: скользящее окно с двумя инвариантами — различных
// символов не более k, каждый символ не более 2 раз. Правый указатель
// всегда идёт вперёд; при нарушении любого инварианта левый сдвигается,
// пока окно снова валидно. Каждый символ входит и выходит из окна один
// раз — O(n) времени, O(distinct) памяти.
func longestSubstring(s string, k int) (int, string) {
	if k <= 0 {
		return 0, ""
	}
	freq := make(map[byte]int)
	bestLen, bestStart := 0, 0
	left := 0

	for right := 0; right < len(s); right++ {
		freq[s[right]]++
		for len(freq) > k || freq[s[right]] > 2 {
			freq[s[left]]--
			if freq[s[left]] == 0 {
				delete(freq, s[left])
			}
			left++
		}
		if right-left+1 > bestLen {
			bestLen = right - left + 1
			bestStart = left
		}
	}
	return bestLen, s[bestStart : bestStart+bestLen]
}

func main() {
	fmt.Println(longestSubstring("abacccbbad", 3)) // 5 "abacc" (или "ccbba")
	fmt.Println(longestSubstring("aaaa", 1))       // 2 "aa" — ограничение «не более 2 раз»
	fmt.Println(longestSubstring("abcdef", 2))     // 2 "ab" — ограничение k различных
	fmt.Println(longestSubstring("", 3))           // 0 ""
	fmt.Println(longestSubstring("aabbcc", 3))     // 6 "aabbcc"
	fmt.Println(longestSubstring("aabbcc", 2))     // 4 "aabb"
}
