package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"sync/atomic"
)

func main() {
	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://ёёёё",
	}

	externalSignal := getExternalSignal()
	sign := atomic.Int32{}
	go func() {
		<-externalSignal
		sign.Store(1)
	}()
	ch := make(chan string)
	wg := sync.WaitGroup{}

	wg.Add(len(urls))
	for i := 0; i < len(urls); i++ {
		i := i
		go func() {
			if sign.Load() == 1 {
				return
			}

			resp, err := http.Get(urls[i])
			if err != nil {
				slog.Error("err", "http get", err)
			}

			if resp.StatusCode == 200 {
				ch <- fmt.Sprintf("адрес %s - ok", urls[i])
			} else {
				ch <- "not ok"
			}

			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for msg := range ch {
		fmt.Println(msg)
	}
}

func getExternalSignal() chan struct{} {
	return make(chan struct{})
}
