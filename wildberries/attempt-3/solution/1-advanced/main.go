package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var cacheStore sync.Map

const (
	timeout       = time.Second
	maxAttempts   = 3
	repairWaiting = 30 * time.Second
)

type result struct {
	data interface{}
	err  error
}

type Processor struct {
	longCounter    *atomic.Int32
	mu             *sync.Mutex
	firstErrorTime time.Time
}

func NewProcessor() *Processor {
	return &Processor{
		longCounter: &atomic.Int32{},
		mu:          &sync.Mutex{},
	}
}

// GetData возвращает данные из getter, страхуя кешом
// в случае ошибки. Проблема: getter может отвечать очень долго.
// Задача: в случае ответа дольше timeout отдавать кеш,
// и правильно обработать кейс, когда getter всегда отвечает долго.
func (p *Processor) GetData(key string, getter func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	ch := make(chan result)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go func() {
		data, err := getter(ctx)
		ch <- result{
			data: data,
			err:  err,
		}
	}()

	if p.longCounter.Load() >= maxAttempts {
		p.mu.Lock()
		if time.Since(p.firstErrorTime) > repairWaiting && !p.firstErrorTime.IsZero() {
			p.longCounter.Store(0)
			p.firstErrorTime = time.Time{}
		}
		p.mu.Unlock()
	}

	select {
	case res := <-ch:
		if res.err != nil {
			fmt.Printf("Getter result err: %s", res.err)

			if data, ok := cacheStore.Load(key); ok {
				return data, nil
			}

			return nil, res.err
		}

		p.mu.Lock()
		p.longCounter.Store(0)
		p.firstErrorTime = time.Time{}
		p.mu.Unlock()

		cacheStore.Store(key, res.data)

		return res.data, nil
	case <-time.After(timeout):
		fmt.Printf("Getter result timeout")

		p.longCounter.Add(1)

		if p.longCounter.Load() == maxAttempts {
			p.mu.Lock()
			p.firstErrorTime = time.Now()
			p.mu.Unlock()
			return nil, fmt.Errorf("too many timeouts")
		}

		if p.longCounter.Load() > maxAttempts {
			return nil, fmt.Errorf("too many timeouts")
		}

		if data, ok := cacheStore.Load(key); ok {
			return data, nil
		}

		return nil, fmt.Errorf("getter timed out and no cache available")
	}
}

func main() {
	p := NewProcessor()
	data, err := p.GetData("key", func(ctx context.Context) (interface{}, error) {
		return "some data", nil
	})

	fmt.Println(data, err)
}
