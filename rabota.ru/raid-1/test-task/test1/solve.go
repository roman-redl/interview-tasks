package main

// Решение задачи QueueAPI (rabota.ru, тест 1).
//
// Проблема 1 — битые сообщения. ctx.Request.Body() у fasthttp возвращает
// слайс из внутреннего буфера, который переиспользуется после ответа:
// горутина, уже отпустившая запрос, читает поверх его тела тело следующего.
// Отсюда «сообщения из кусков разных сообщений». Лечение — копировать
// body до того, как управление вернулось в fasthttp.
//
// Проблема 2 — смерть под нагрузкой. `go putToQueue(...)` на каждый
// запрос — неограниченные горутины и память; сервис должен «пусть
// тормозить, но выживать, затормаживая и потребителей»: ограниченный пул
// воркеров, конечная очередь между handler и воркерами, а при заполненной
// очереди — явный отказ с Retry-After и глубиной очереди в заголовке.
// Потребитель сам решает: повторить позже (429/503) или слать дальше.

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

const (
	addr        = ":8088"
	queueHeader = "X-Queue-Depth" // сколько задач уже ждёт воркеров
)

// putToQueue имитирует сетевую отправку (как в условии, в секундах).
func putToQueue(body []byte) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Second)
	log.Println(body)
}

// queue: tasks — задачи на отправку, depth — текущая заполненность.
type queue struct {
	tasks chan []byte
}

func newQueue(workers, backlog int) *queue {
	q := &queue{tasks: make(chan []byte, backlog)}
	for i := 0; i < workers; i++ {
		go func() {
			for body := range q.tasks {
				putToQueue(body)
			}
		}()
	}
	return q
}

// tryEnqueue: nil, false — очередь полна, клиенту пора притормозить.
func (q *queue) tryEnqueue(body []byte) bool {
	select {
	case q.tasks <- body:
		return true
	default:
		return false
	}
}

func envInt(name string, def int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

func main() {
	workers := envInt("QUEUEAPI_WORKERS", runtime.NumCPU())
	backlog := envInt("QUEUEAPI_BACKLOG", 1024)
	q := newQueue(workers, backlog)

	r := router.New()

	r.PUT("/queue", func(ctx *fasthttp.RequestCtx) {
		// копия до возврата в fasthttp: буфер запроса переиспользуется
		body := append([]byte{}, ctx.Request.Body()...)

		ctx.Response.Header.Set(queueHeader, strconv.Itoa(len(q.tasks)))
		if !q.tryEnqueue(body) {
			// очередь переполнена: не копим горутины, отдаём backpressure —
			// потребитель решает сам, ждать и повторить или отвалиться
			ctx.Response.Header.Set("Retry-After", "1")
			ctx.Error(
				fmt.Sprintf("queue is full: %d pending, try later", backlog),
				fasthttp.StatusTooManyRequests,
			)
			return
		}
		ctx.Response.SetStatusCode(fasthttp.StatusAccepted)
	})

	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("QueueAPI on %s: workers=%d backlog=%d", addr, workers, backlog)
	log.Fatal(fasthttp.Serve(ln, r.Handler))
}
