# Решение 7 — ревью сервиса отправки уведомлений

Проблемы по требованиям, от критичных к стилистическим.

**Требование «остановить остальные при ошибке» — не выполнено.** Запись
`resultErr = err` ничего не останавливает: остальные горутины продолжают
слать запросы. Нужна отмена: `context.WithCancel` по первой ошибке (или
`errgroup.WithContext`), проверка `ctx.Err()` в воркерах и в подачe задач.

**Data race на `resultErr`.** Параллельные горутины пишут одну переменную
без синхронизации — undefined behavior, поймается `-race`. Фикс:
`sync.Once` для первого значения (или канал ошибок).

**Data race на `sentCount`.** `c.sentCount++` и `SentCount()` без
синхронизации при конкурентных Send. Фикс: `atomic.Int64` (Add/Load).

**Контекст не используется.** `http.NewRequest` вместо
`http.NewRequestWithContext(ctx, ...)` — отмена не доезжает до запроса;
воркеры не проверяют `ctx.Done()`. Требование «корректно обрабатывать
отмену» не выполнено: отменённый Send продолжит слать всё.

**Утечка коннектов: `resp.Body` не закрывается и не дренируется.** Каждый
ответ держит соединение; при массовом отправлении пул исчерпается. Фикс:
`defer` закрытие + `io.Copy(io.Discard, resp.Body)` для keep-alive.

**Игнорируются ошибки Marshal/NewRequest** (`body, _ :=`,
`req, _ :=`). Для «все ресурсы корректно освобождаются» это прямой баг:
req == nil → паника на следующей строке. Обрабатывать и возвращать.

**Горутина на каждое сообщение.** Требование «не более maxWorkers
HTTP-запросов» семафор держит, но на длинном входе создаются миллионы
горутин, ждущих слот — память и планировщик. Каноничное решение — пул из
maxWorkers воркеров, читающих канал задач.

**`time.Sleep(time.Second)` в конце работы** — назначение неясно (похоже
на «rate limit», но блокирует воркера на секунду после каждого запроса,
срезая пропускную способность в разы). Убрать; если нужен rate limit —
rate.Limiter на отправку.

**Замыкание захватывает `msg` цикла** — в Go до 1.22 переменная range
одна на все итерации: воркеры могут отправить одно сообщение несколько
раз. Даже на новых версиях параметр `go func(msg Message)` — дешёвая
страховка.

**Нет валидации `maxWorkers <= 0`** — `make(chan struct{}, 0)` /
отрицательный: паника или вечная блокировка. `len(messages) == 0` —
быстрый выход.

**Мелочи**: 4xx не считаются ошибкой, хотя это неудачная отправка
(требование счётчика «успешно отправленных»); `fmt.Errorf("server
error")` без кода — не диагностировать; `http.Client` без Timeout
переложен — таймаут есть, хорошо (единственное, что сделано верно).

Исправленная версия (скомпилирована, прогнана с `-race`, поведение
проверено: 20 сообщений, каждое 4-е отвечает 500 — Send возвращает ошибку,
счётчик = успешные до отмены):

```go
package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Message struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

type Client struct {
	httpClient *http.Client
	url        string

	sentCount atomic.Int64
}

func NewClient(url string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		url:        url,
	}
}

func (c *Client) SentCount() int64 {
	return c.sentCount.Load()
}

func (c *Client) Send(ctx context.Context, messages []Message, maxWorkers int) error {
	if maxWorkers <= 0 {
		return fmt.Errorf("maxWorkers must be positive, got %d", maxWorkers)
	}
	if len(messages) == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan int)
	workersDone := make(chan struct{})

	var (
		firstErr error
		once     sync.Once // единственный писатель firstErr; чтение — после workersDone
	)
	reportErr := func(err error) {
		cancel() // остановить остальных: их ctx отменён
		once.Do(func() { firstErr = err })
	}

	// пул воркеров фиксированного размера: не горутина на сообщение
	go func() {
		defer close(workersDone)
		var wg sync.WaitGroup
		wg.Add(maxWorkers)
		for w := 0; w < maxWorkers; w++ {
			go func() {
				defer wg.Done()
				for i := range jobs {
					if ctx.Err() != nil {
						return
					}
					if err := c.sendOne(ctx, messages[i]); err != nil {
						reportErr(err)
						continue
					}
					c.sentCount.Add(1)
				}
			}()
		}
		wg.Wait()
	}()

feed:
	for i := range messages {
		select {
		case jobs <- i:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	<-workersDone

	if firstErr != nil {
		return firstErr
	}
	return ctx.Err()
}

func (c *Client) sendOne(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body) // дренаж для keep-alive
		resp.Body.Close()
	}()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("server error: status %d", resp.StatusCode)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("client error: status %d", resp.StatusCode)
	}
	return nil
}
```

Альтернатива короче — `golang.org/x/sync/errgroup`:
`g, gctx := errgroup.WithContext(ctx)`, `g.SetLimit(maxWorkers)` —
пул, отмена по первой ошибке и сбор ошибки даёт один вызов вместо
ручной сборки; показательное дополнение на ревью: «в проде взял бы
errgroup».
