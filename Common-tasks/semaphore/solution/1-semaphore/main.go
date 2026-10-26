package main

import (
	"fmt"
	"sync"
	"time"
)

type Semaphore struct {
	ch chan struct{}
}

func NewSemaphore(weight int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, weight),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

func (s *Semaphore) Release() {
	<-s.ch
}

func main() {
	sem := NewSemaphore(2)
	var wg sync.WaitGroup

	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem.Acquire()
			fmt.Printf("Горутина %d захватила ресурс\n", id)
			time.Sleep(time.Second * 3)
			fmt.Printf("Горутина %d освобождает ресурс\n", id)
			sem.Release()
		}(i)
	}

	wg.Wait()
}
