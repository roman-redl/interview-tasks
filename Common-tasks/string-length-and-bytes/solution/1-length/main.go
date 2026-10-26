package main

import "fmt"

func main() {
	str := "Привет, мир!"

	length := len(str)
	fmt.Println("Длина строки в символах:", length)

	length = len([]byte(str))
	fmt.Println("Длина строки в байтах:", length)
}
