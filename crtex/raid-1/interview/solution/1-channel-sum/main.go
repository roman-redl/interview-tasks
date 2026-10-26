package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(v int) {
			defer wg.Done()
			ch <- v
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var sum int
	for v := range ch {
		sum += v
	}

	fmt.Printf("result: %d\n", sum)
}
