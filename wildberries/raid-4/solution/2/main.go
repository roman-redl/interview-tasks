package main

import (
	"fmt"
)

func main() {
	// Нужно вывести символы английского алфавита
	// - каждая третья буква должна заменяться и вместо нее выводиться индекс
	// - каждая пятая буква должна заменяться и вместо нее выводится буква русского алфавита согласно индексу
	// - каждая пятнадцатая - выводить 15. // 1 раз
	alphabet := "abcdefgihgklmnoprstyfxz" //...
	russianAlphabet := []rune("цукенгшщзхъфывапролджэячсмитьбю")

	for i := 1; i < len(alphabet); i++ {
		if i%15 == 0 {
			fmt.Println(15)

			continue
		}
		if i%3 == 0 {
			fmt.Println(i)

			continue
		}
		if i%5 == 0 {
			fmt.Println(string(russianAlphabet[i-1]))

			continue
		}

		fmt.Println(string(alphabet[i-1]))
	}
}
