package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	count := 10
	limit := 5
	wg := sync.WaitGroup{}

	ch := make(chan int, limit)
	wg.Add(count)

	for i := 1; i <= count; i++ {
		go func() {
			defer wg.Done()

			ch <- i
			printNumber(i) // Написано специально для версии 1.22, в 1.21 не работает из-за проскальзывания чисел😁
			<-ch
		}()
	}

	wg.Wait()
	close(ch)
}

func printNumber(n int) {
	time.Sleep(1 * time.Second)
	fmt.Println(n)
}
