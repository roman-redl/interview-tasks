package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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
	invokeCounter *atomic.Int32
	errorsCounter int
	block         bool
}

// addr содержит ip:port конкретного экземпляра
func NewBackend(addr string) *BackendImpl {
	return &BackendImpl{}
}

type Balancer struct {
	backends           []*BackendWithCounter
	minBackend         *BackendWithCounter
	currentBackCounter *atomic.Int32
	mu                 *sync.Mutex
}

var _ Backend = &Balancer{}

// addrs содержат адреса всех балансируемых экземпляров
func NewBalancer(addrs []string) *Balancer {
	backends := make([]*BackendWithCounter, len(addrs))

	for i, addr := range addrs {
		backends[i] = &BackendWithCounter{
			invokeCounter: &atomic.Int32{},
			b:             NewBackend(addr),
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

	for i := 0; i < len(b.backends); i++ {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		b.mu.Lock()
		minInvokes := int32(-1)
		if b.minBackend == nil || b.minBackend.block {
			for _, back := range b.backends {
				if !back.block && (minInvokes == -1 || back.invokeCounter.Load() < minInvokes) {
					b.minBackend = back
					minInvokes = back.invokeCounter.Load()
				}
			}
		}
		if b.minBackend == nil {
			b.mu.Unlock()
			break
		}
		selected := b.minBackend
		b.mu.Unlock()

		selected.invokeCounter.Add(1)
		resp, err := selected.b.Invoke(ctx, req)
		selected.invokeCounter.Add(-1)
		if err != nil {
			b.mu.Lock()
			selected.errorsCounter++
			if selected.errorsCounter == 1 {
				selected.t = time.Now()
			}
			if selected.errorsCounter >= n && time.Since(selected.t) <= x {
				selected.block = true
				b.minBackend = nil
			}
			b.mu.Unlock()

			continue
		}

		b.mu.Lock()
		selected.errorsCounter = 0
		selected.t = time.Time{}
		b.minBackend = nil
		b.mu.Unlock()

		return resp, nil
	}

	return nil, fmt.Errorf("all backends are blocked")
}

// Решение — библиотечный код; main нужен, чтобы пакет компилировался и запускался сам по себе.
func main() {}
