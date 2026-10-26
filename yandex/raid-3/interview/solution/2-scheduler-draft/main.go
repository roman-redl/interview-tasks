package main

import (
	"sync"
)

/*
Реализовать шедулер, который будет

- быстро принимать входящие задания,
- отдавать результат задания по запросу
- если результата нет, возвращать соответствующий статус
*/

/*
Процессор - сервис, который выполняет долгую ресурсоемкую операцию. Он уже реализован
Возвращаемые им ошибки связаны исключительно с некорректными входными данными, устойчивы и перезапрашивать их бесполезно
*/
type Processor interface {
	Process([]byte) ([]byte, error)
}

/*
Считаем, что у нас есть генератор уникальных идентификаторов (например, UUID).
Генератор гарантирует, что коллизий не будет
*/
type ID string

func NewID() ID {
	return "uuid"
}

type job struct {
	status string
	result []byte
	err    error
}

/*
Шедулер принимает от клиентов задачи, ставит в очередь и запускает в процессинг.
Обеспечивает, что единовременно запущено не более threads методов process
Дает возможность проверить статус задачи и получить результат
не блокирует методы своего публичного интерфейса на процессинг
*/
type Scheduler struct {
	processor Processor
	mu        sync.Mutex
	result    map[ID]job
	sem       chan struct{}
}

func NewScheduler(
	prc Processor,
	threads int,
) *Scheduler {
	return &Scheduler{
		processor: prc,
		result:    make(map[ID]job),
		sem:       make(chan struct{}, threads),
	}
}

func (s *Scheduler) Queue(request []byte) ID {
	id := NewID()

	/*
		s.mu.Lock()
		job := s.result[id] // тут не совсем правильно было на собесе, нужно было создать новый job
		job.status = "queued"
		s.result[id] = job
		s.mu.Unlock()
	*/

	s.mu.Lock()
	s.result[id] = job{status: "pending"}
	s.mu.Unlock()

	go func() {
		// Ограничиваем количество горутин в блоке процессинга
		s.sem <- struct{}{}
		defer func() { <-s.sem }()

		s.mu.Lock()
		s.result[id] = job{status: "processing"}
		s.mu.Unlock()

		result, err := s.processor.Process(request)
		if err != nil {
			s.mu.Lock()
			s.result[id] = job{status: "error", err: err}
			s.mu.Unlock()
		} else {
			s.mu.Lock()
			s.result[id] = job{status: "completed", result: result}
			s.mu.Unlock()
		}
	}()

	return id
}

func (s *Scheduler) Status(id ID) (status string, result []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock() // забыл сделать unlock на собесе

	res, ok := s.result[id]
	if !ok {
		return "not_found", nil, nil
	}

	return res.status, res.result, res.err
}

// Решение — библиотечный код; main нужен, чтобы пакет компилировался и запускался сам по себе.
func main() {}
