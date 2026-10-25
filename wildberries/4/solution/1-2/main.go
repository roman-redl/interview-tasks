package main

import (
	"fmt"
)

func add(a int) func(b int) int {
	return func(b int) int {
		return a + b
	}
}

func main() {
	// 1 Может ли работать такая функция
	d := add(3)(5)
	fmt.Println(d) // 8

	// 2 Нужно вывести символы английского алфавита
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
