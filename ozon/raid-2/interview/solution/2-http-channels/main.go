package main

import (
	"fmt"
	"net/http"
)

type resp struct {
	code int
}

func Get(url string) (resp, error) {
	r, err := http.Get(url)
	if err != nil {
		return resp{}, err
	}
	defer r.Body.Close()
	return resp{code: r.StatusCode}, nil
}

type result struct {
	url string
	ok  bool
}

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

	// Каждый URL запрашивается в своей горутине; горутина отправляет
	// результат в канал; печать происходит только в основном потоке.
	// Буфер len(urls): горутины не блокируются на отправке, main читает
	// столько же ответов, сколько было задач — количество синхронизирует
	// цикл приёма, отдельный признак конца не нужен.
	results := make(chan result, len(urls))
	for _, u := range urls {
		go func(u string) {
			r, err := Get(u)
			results <- result{url: u, ok: err == nil && r.code == http.StatusOK}
		}(u)
	}

	for range urls {
		r := <-results
		status := "not ok"
		if r.ok {
			status = "ok"
		}
		fmt.Printf("%s - %s\n", r.url, status)
	}
	// Порядок печатей недетерминирован — кто раньше ответил, тот раньше
	// попал в канал. Нужен порядок входа — отправлять в канал пару (индекс,
	// результат) и собирать в слайс до печати.
}
