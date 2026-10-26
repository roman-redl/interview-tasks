package main

import (
	"fmt"
	"sync"
)

// Шардированный кэш: ключ хэшируется (FNV-1a), по хэшу выбирается один из
// N независимых шардов, каждый под своим sync.RWMutex. Конкуренция за
// глобальную блокировку падает примерно в N раз; Get и Set разных ключей
// почти всегда не встречаются.
const shardCount = 16

type shard struct {
	mu sync.RWMutex
	m  map[string]string
}

type shardedCache struct {
	shards [shardCount]*shard
}

func NewSharded() *shardedCache {
	c := &shardedCache{}
	for i := range c.shards {
		c.shards[i] = &shard{m: make(map[string]string)}
	}
	return c
}

func (c *shardedCache) shardFor(k string) *shard {
	h := uint32(2166136261) // FNV-1a offset basis
	for i := 0; i < len(k); i++ {
		h ^= uint32(k[i])
		h *= 16777619
	}
	return c.shards[h%shardCount]
}

func (c *shardedCache) Set(k, v string) {
	s := c.shardFor(k)
	s.mu.Lock()
	s.m[k] = v
	s.mu.Unlock()
}

func (c *shardedCache) Get(k string) (string, bool) {
	s := c.shardFor(k)
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[k]
	return v, ok
}

func main() {
	c := NewSharded()

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
