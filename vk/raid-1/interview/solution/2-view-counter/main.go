package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// ViewCounter — шардированный счётчик: конкурентный инкремент разных
// ключей не встречается на одном локе; значения — atomic, чтение без
// блокировки. Персист — периодический snapshot в файл (tmp + rename),
// окно потери при сбое — один интервал flush.
type ViewCounter struct {
	shardCount int
	shards     []shard
	// дельта с прошлого snapshot: пишем только изменённое
	dirty sync.Map

	snapshotPath string
	stop         chan struct{}
	done         chan struct{}
}

type shard struct {
	mu sync.RWMutex
	m  map[string]*atomic.Int64
}

func NewViewCounter(shards int, snapshotPath string) *ViewCounter {
	vc := &ViewCounter{
		shardCount:   shards,
		shards:       make([]shard, shards),
		snapshotPath: snapshotPath,
		stop:         make(chan struct{}),
		done:         make(chan struct{}),
	}
	for i := range vc.shards {
		vc.shards[i].m = make(map[string]*atomic.Int64)
	}
	vc.loadSnapshot()
	go vc.flushLoop(2 * time.Second)
	return vc
}

func (vc *ViewCounter) shardFor(key string) int {
	h := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		h ^= uint32(key[i])
		h *= 16777619
	}
	return int(h) % vc.shardCount
}

func (vc *ViewCounter) counter(pageID string) *atomic.Int64 {
	s := &vc.shards[vc.shardFor(pageID)]
	s.mu.RLock()
	c := s.m[pageID]
	s.mu.RUnlock()
	if c != nil {
		return c
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if c = s.m[pageID]; c == nil {
		c = &atomic.Int64{}
		s.m[pageID] = c
	}
	return c
}

func (vc *ViewCounter) Inc(pageID string) {
	vc.counter(pageID).Add(1)
	vc.dirty.Store(pageID, true)
}

func (vc *ViewCounter) Get(pageID string) int64 {
	s := &vc.shards[vc.shardFor(pageID)]
	s.mu.RLock()
	c := s.m[pageID]
	s.mu.RUnlock()
	if c == nil {
		return 0
	}
	return c.Load()
}

type snapshotEntry struct {
	Key   string
	Value int64
}

func (vc *ViewCounter) loadSnapshot() {
	f, err := os.Open(vc.snapshotPath)
	if err != nil {
		return
	}
	defer f.Close()
	var entries []snapshotEntry
	if err := gob.NewDecoder(f).Decode(&entries); err != nil {
		log.Println("snapshot decode:", err)
		return
	}
	for _, e := range entries {
		s := &vc.shards[vc.shardFor(e.Key)]
		c := &atomic.Int64{}
		c.Store(e.Value)
		s.m[e.Key] = c
	}
	log.Printf("restored %d counters from snapshot", len(entries))
}

func (vc *ViewCounter) flush() {
	var entries []snapshotEntry
	vc.dirty.Range(func(key, _ any) bool {
		k := key.(string)
		entries = append(entries, snapshotEntry{Key: k, Value: vc.Get(k)})
		return true
	})
	if len(entries) == 0 {
		return
	}

	f, err := os.CreateTemp(filepath.Dir(vc.snapshotPath), ".snap-*")
	if err != nil {
		log.Println("snapshot create:", err)
		return
	}
	if err := gob.NewEncoder(f).Encode(entries); err != nil {
		f.Close()
		log.Println("snapshot encode:", err)
		return
	}
	if err := f.Close(); err != nil {
		log.Println("snapshot close:", err)
		return
	}
	// атомарная замена: файл либо целиком старый, либо целиком новый
	if err := os.Rename(f.Name(), vc.snapshotPath); err != nil {
		log.Println("snapshot rename:", err)
		return
	}
	for _, e := range entries {
		vc.dirty.Delete(e.Key)
	}
}

func (vc *ViewCounter) flushLoop(interval time.Duration) {
	defer close(vc.done)
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			vc.flush()
		case <-vc.stop:
			vc.flush() // финальный snapshot на shutdown
			return
		}
	}
}

func (vc *ViewCounter) Close() {
	close(vc.stop)
	<-vc.done
}

func main() {
	vc := NewViewCounter(64, "/tmp/viewcounter.gob")

	mux := http.NewServeMux()
	mux.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		pageID := strings.TrimPrefix(r.URL.Path, "/view/")
		switch r.Method {
		case http.MethodPost:
			vc.Inc(pageID)
			w.WriteHeader(http.StatusAccepted)
		case http.MethodGet:
			fmt.Fprintf(w, "%d\n", vc.Get(pageID))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	srv := &http.Server{Addr: ":8089", Handler: mux}
	go func() {
		log.Println("view counter on :8089")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig // graceful: додать начатые запросы и снять финальный snapshot
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	vc.Close()
}
