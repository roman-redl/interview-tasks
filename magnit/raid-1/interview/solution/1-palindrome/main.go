package main

import "fmt"

// isPalindrome: два указателя с концов к середине; O(n) времени,
// O(1) памяти. Полное условие задачи не сохранилось — базовая
// формулировка: строка палиндром, если читается одинаково в обе стороны;
// регистр и небуквенные символы игнорируются (ASCII; для Unicode —
// работать рунами и unicode.IsLetter/IsDigit).
func isPalindrome(s string) bool {
	l, r := 0, len(s)-1
	for l < r {
		for l < r && !isAlnum(s[l]) {
			l++
		}
		for l < r && !isAlnum(s[r]) {
			r--
		}
		if lower(s[l]) != lower(s[r]) {
			return false
		}
		l++
		r--
	}
	return true
}

func isAlnum(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

func main() {
	fmt.Println(isPalindrome("racecar"))    // true
	fmt.Println(isPalindrome("Race, car"))  // true — регистр и пунктуация игнорируются
	fmt.Println(isPalindrome("race a car")) // false
	fmt.Println(isPalindrome(""))           // true — пустая строка палиндром
	fmt.Println(isPalindrome("1 2 1"))      // true
}
