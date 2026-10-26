package main

import "fmt"

func main() {
	str := "Hello, world!"

	newStr := replaceChar(str, 'l', 'X')
	fmt.Println(newStr)
}

func replaceChar(input string, oldChar rune, newChar rune) string {
	bytes := []byte(input)

	for i, b := range bytes {
		if rune(b) == oldChar {
			bytes[i] = byte(newChar)
		}
	}

	return string(bytes)
}
