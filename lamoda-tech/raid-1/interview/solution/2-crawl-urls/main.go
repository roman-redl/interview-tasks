package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// crawl обходит urls параллельно — не более k запросов одновременно —
// и возвращает коды ответов в исходном порядке.
//
// Фиксированный пул из k горутин-воркеров, а не горутина на URL: большой
// вход не создаёт шторма горутин на планировщике и GC. Канал задач с
// буфером 2*k: небуферизированный канал заставляет поставщика и воркеров
// «здороваться» на каждой задаче (старт-стоп), буфер длиной len(urls)
// дорог по памяти; небольшой буфер, кратный k, сглаживает разницу скоростей
// поставщика и воркеров, не занимая память под весь вход.
func crawl(urls []string, k int) []int {
	type job struct {
		idx int
		url string
	}

	jobs := make(chan job, 2*k)
	out := make([]int, len(urls))

	client := &http.Client{Timeout: 5 * time.Second}

	var wg sync.WaitGroup
	wg.Add(k)
	for w := 0; w < k; w++ {
		go func() {
			defer wg.Done()
			for j := range jobs {
				resp, err := client.Get(j.url)
				switch {
				case err != nil:
					out[j.idx] = 0 // договорённость: сетевая ошибка => 0
				default:
					out[j.idx] = resp.StatusCode
					resp.Body.Close()
				}
			}
		}()
	}

	for i, u := range urls {
		jobs <- job{idx: i, url: u}
	}
	close(jobs)
	wg.Wait()

	return out
}

func main() {
	// Вместо реальных URL — локальный httptest-сервер: решение проверяемо
	// офлайн. Обход 20 адресов с ограничением 5 показывает сохранение
	// порядка и различие кодов ответов.
	codes := map[string]int{"/ok": 200, "/moved": 302, "/err": 500}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		code := codes[r.URL.Path]
		if code == 0 {
			code = 404
		}
		w.WriteHeader(code)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	var urls []string
	want := []int{}
	paths := []string{"/ok", "/moved", "/err", "/nope"}
	wantCodes := []int{200, 302, 500, 404}
	for i := 0; i < 20; i++ {
		urls = append(urls, srv.URL+paths[i%4])
		want = append(want, wantCodes[i%4])
	}

	res := crawl(urls, 5)
	status := "OK"
	for i := range res {
		if res[i] != want[i] {
			status = "FAIL"
		}
	}
	fmt.Println(status, res)
}
