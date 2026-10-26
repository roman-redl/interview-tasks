# UPSL — заход 1, скрининг на HR-созвоне

- Даты: 2026-06-01
- Формат: онлайн, HR-созвон: блиц-теория + ревью кода в общей live-coding-комнате

## 1. Этапы CI/CD для Go-сервиса

Какие этапы обычно стоят в пайплайне для сервиса на Go? Речь про сборку и деплой приложения в плане CI/CD.

## 2. Контекст в Go

Для чего нужен контекст в Go?

## 3. Контекст в запросе к базе

Что произойдёт, если не прокинуть контекст в запросе к базе?

## 4. Линтеры

С какими линтерами работали, какие основные инструменты использовали?

## 5. Опыт с SQL

Какой опыт работы с SQL? Умеете ли писать и оптимизировать запросы, с какими СУБД чаще всего работали?

## 6. Когда индекс не поможет

Когда индекс не поможет или даже навредит?

## 7. Ревью кода: сервис отправки уведомлений

В live-coding-комнате открыт код. Нужно провести ревью: проговорить как можно больше проблем; то, что правится просто, исправить самому; то, что требует долгого исправления, — оставить комментарием, что нужно изменить.

```go
/*
Разработчику было дано задание написать сервис отправки уведомлений.

Требования:

- Метод Send должен отправлять сообщения параллельно.
- Одновременно может выполняться не более maxWorkers HTTP-запросов.
- Если хотя бы одна отправка завершилась ошибкой — нужно остановить остальные отправки и вернуть ошибку.
- Метод должен корректно обрабатывать отмену контекста.
- Все ресурсы должны корректно освобождаться.
- Метод может вызываться конкурентно из разных goroutine.
- Client должен хранить количество успешно отправленных сообщений.

Нужно сделать code review, найти проблемы и исправить код.
*/

package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	UserID string `json:"user_id"`
	Text   string `json:"text"`
}

type Client struct {
	httpClient *http.Client
	url        string

	sentCount int
}

func NewClient(url string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
		url: url,
	}
}

func (c *Client) SentCount() int {
	return c.sentCount
}

func (c *Client) Send(ctx context.Context, messages []Message, maxWorkers int) error {
	sem := make(chan struct{}, maxWorkers)

	var wg sync.WaitGroup
	var resultErr error

	for _, msg := range messages {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() {
				<-sem
			}()

			body, _ := json.Marshal(msg)

			req, _ := http.NewRequest(
				http.MethodPost,
				c.url,
				bytes.NewReader(body),
			)

			req.Header.Set("Content-Type", "application/json")

			resp, err := c.httpClient.Do(req)
			if err != nil {
				resultErr = err
				return
			}

			if resp.StatusCode >= 500 {
				resultErr = fmt.Errorf("server error")
				return
			}

			c.sentCount++

			time.Sleep(time.Second)
		}()
	}

	wg.Wait()

	return resultErr
}
```
