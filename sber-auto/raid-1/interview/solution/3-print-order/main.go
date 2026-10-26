package main

import (
	"fmt"
	"sync"
)

func printText(data []string) {
	wg := sync.WaitGroup{}
	wg.Add(len(data))

	for _, v := range data {
		go func(v string) {
			fmt.Println(v)

			wg.Done()
		}(v)
	}

	wg.Wait()

	fmt.Println("done!")
}

func main() {
	printText([]string{"one", "two", "three"})
}
