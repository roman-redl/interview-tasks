package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"time"
)

// crawl обходит urls со скоростью не более 5 стартов запросов в секунду
// и без потери пропускной способности на медленных ответах.
//
// Rate limiter на стандартных средствах: time.NewTicker(200 мс) кладёт
// токен в канал ёмкости 1. Ёмкость 1 — принципиальна: пока воркеры заняты
// медленными запросами, токены не копятся; после паузы мы НЕ навёрстываем
// удвоенной скоростью, средняя остаётся 5/с. Воркер берёт токен только
// когда уже держит задачу и готов её выполнить — ограничивается именно
// старт запросов, а не выдача задач.
//
// Пул воркеров (не горутина на ссылку): тысячи горутин на цикле —
// memory/планировщик, неограниченная аллокация; фиксированный пул
// переживает любой размер входа. Медленный URL занимает одного воркера,
// остальные продолжают забирать задачи — пропускная способность не
// проседает, пока есть свободные воркеры (их число — единственный тюнинг
// под latency целей).
//
// Завершение: close(jobs) после раздачи + wg.Wait() — воркеры дописывают
// начатое и выходят по закрытому каналу; токен-тикер живёт до main
// (горутина тикера останавливается defer t.Stop() при завершении
// процесса; в демонстрации он останавливается явным done-каналом).
func crawl(urls []string) time.Duration {
	const rate = 5
	workers := 5

	tokens := make(chan struct{}, 1)
	tickerDone := make(chan struct{})
	go func() {
		t := time.NewTicker(time.Second / rate)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				select {
				case tokens <- struct{}{}:
				default: // предыдущий токен ещё не разобран — не копим
				}
			case <-tickerDone:
				return
			}
		}
	}()

	jobs := make(chan string)
	var wg sync.WaitGroup
	var served atomic.Int64
	start := time.Now()

	wg.Add(workers)
	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for url := range jobs {
				<-tokens // разрешение на старт — берём готовым к работе
				resp, err := http.Get(url)
				if err == nil {
					resp.Body.Close()
				}
				fmt.Printf("done #%d %s %v\n", served.Add(1), url, respStatus(resp, err))
			}
		}()
	}

	for _, u := range urls {
		jobs <- u
	}
	close(jobs)
	wg.Wait()
	close(tickerDone)

	return time.Since(start)
}

func respStatus(resp *http.Response, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	return resp.Status
}

func main() {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1)%3 == 0 {
			time.Sleep(300 * time.Millisecond) // каждый третий ответ медленный
		}
	}))
	defer srv.Close()

	urls := make([]string, 10)
	for i := range urls {
		urls[i] = srv.URL
	}

	elapsed := crawl(urls)
	fmt.Println("elapsed:", elapsed.Round(time.Millisecond))
	// 10 запросов при 5/с ≈ 1.8–2.2 с; медленные ответы не растянули обход:
	// свободные воркеры подхватывали задачи, пока один ждал.
}
