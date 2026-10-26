package main

import "fmt"

func main() {
	s := "Привет."
	// s == "Приветы"

	// база
	newStr := ""
	for _, r := range s {
		if r == '.' {
			newStr += "ы"
		} else {
			newStr += string(r)
		}
	}

	fmt.Println(newStr)

	// Оптимальное решение
	result := []rune(s)
	fmt.Println(string(result[:len(result)-1]) + "ы") // либо через WriteByte

}
