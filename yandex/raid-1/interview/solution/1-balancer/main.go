package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/*
Есть приложение с микросервисной архитектурой.
Микросервис можно абстрагировать с помощью интерфейса Backend.
Для доступа к одному экземпляру микросервиса можно использовать
тип BackendImpl, который уже реализован.

Для каждого микросервиса есть несколько десятков запущенных
экземпляров, каждый из которых доступен по своему адресу addr.
Однако отдельные экземпляры микросервиса ненадежны:
они могут падать, быть недоступными либо перегруженными.
Поэтому вам нужно реализовать тип Balancer, который также реализует
интерфейс Backend и осуществляет client-side балансировку нагрузки
между экземплярами микросервиса.

1. Если в течении x времени было n ошибок подряд, то блокирует на y минут, используем round robin.
2. Усложнение: модернизировать и использовать для запросов наименее нагруженные бэкенды.

Решить за час
*/

const (
	x     = 1 * time.Minute
	n int = 10
	y     = 5 * time.Minute
)

type Request interface{}

type Response interface{}

type Backend interface {
	Invoke(ctx context.Context, req Request) (Response, error)
}

type BackendImpl struct{}

func (b BackendImpl) Invoke(ctx context.Context, req Request) (Response, error) {
	//TODO implement me
	panic("implement me")
}

var _ Backend = &BackendImpl{}

type BackendWithCounter struct {
	b             Backend
	t             time.Time
	errorsCounter int
	block         bool
}

// addr содержит ip:port конкретного экземпляра
func NewBackend(addr string) *BackendImpl {
	return &BackendImpl{}
}

type Balancer struct {
	backends           []*BackendWithCounter
	currentBackCounter int
	mu                 *sync.Mutex
}

var _ Backend = &Balancer{}

// addrs содержат адреса всех балансируемых экземпляров
func NewBalancer(addrs []string) *Balancer {
	backends := make([]*BackendWithCounter, len(addrs))

	for i, addr := range addrs {
		backends[i] = &BackendWithCounter{
			b: NewBackend(addr),
		}
	}

	return &Balancer{
		backends: backends,
		mu:       &sync.Mutex{},
	}
}

func (b *Balancer) Invoke(ctx context.Context, req Request) (Response, error) {
	b.mu.Lock()
	for _, back := range b.backends {
		if time.Since(back.t) > y && !back.t.IsZero() {
			back.block = false
			back.errorsCounter = 0
			back.t = time.Time{}
		}
	}
	b.mu.Unlock()

	for range len(b.backends) {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		b.mu.Lock()
		back := b.backends[b.currentBackCounter%len(b.backends)]
		b.currentBackCounter++
		if back.block {
			continue
		}
		b.mu.Unlock()

		resp, err := back.b.Invoke(ctx, req)
		if err != nil {
			b.mu.Lock()
			back.errorsCounter++
			if back.errorsCounter == 1 {
				back.t = time.Now()
			}

			if back.errorsCounter >= n && time.Since(back.t) <= x {
				back.block = true
				b.mu.Unlock()

				continue
			} else {
				b.mu.Unlock()
				return nil, err
			}
		}

		b.mu.Lock()
		back.errorsCounter = 0
		back.t = time.Time{}
		b.mu.Unlock()

		return resp, nil
	}

	return nil, fmt.Errorf("all backends are blocked")
}

// Решение — библиотечный код; main нужен, чтобы пакет компилировался и запускался сам по себе.
func main() {}
