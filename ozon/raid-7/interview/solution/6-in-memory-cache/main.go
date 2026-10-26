package main

import (
	"fmt"
	"sync"
)

// Простая потокобезопасная реализация: map под одним sync.RWMutex.
// Параллельные Get идут под RLock; Set — под Lock. При росте нагрузки
// одна блокировка становится узким местом — шардированный вариант в
// 6-in-memory-cache-sharded.
type Cache interface {
	Set(k, v string)
	Get(k string) (v string, ok bool)
}

type cache struct {
	mu sync.RWMutex
	m  map[string]string
}

func New() Cache {
	return &cache{m: make(map[string]string)}
}

func (c *cache) Set(k, v string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k] = v
}

func (c *cache) Get(k string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.m[k]
	return v, ok
}

func main() {
	c := New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func(i int) { defer wg.Done(); c.Set(fmt.Sprintf("k%d", i), "v") }(i)
		go func(i int) { defer wg.Done(); c.Get(fmt.Sprintf("k%d", i)) }(i)
	}
	wg.Wait()

	v, ok := c.Get("k42")
	fmt.Println("k42:", v, ok) // k42: v true
}
