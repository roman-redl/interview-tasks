# М.Видео — заход 2, техническое интервью

- Даты: 2026-04-15
- Формат: онлайн, live-coding + разбор опыта, теория (Kafka, Temporal, PostgreSQL)

## 1. Пакетная обработка батчами по 50

В онлайн-редакторе дан заготовленный код. Реализовать метод `ProcessBatches`: он принимает слайс строк и обрабатывает его батчами (группами) по 50 элементов — для каждого батча вызвать сохранение и отправку уведомления.

```go
package main

import "context"

// Необходимо реализовать функцию ProcessBatches, которая принимает слайс строк
// и обрабатывает его батчами (группами) по 50 штук

type saver interface {
	Save(data []string) error
}

type notifier interface {
	Notify(data []string) error
}

type Service struct {
	saver    saver
	notifyer notifier
}

func New(saver saver, notifier notifier) *Service {
	return &Service{
		saver:    saver,
		notifyer: notifier,
	}
}

// ProcessBatches разбивает входные данные на бати указанного размера
func (s *Service)ProcessBatches(ctx context.Context, items []string) error {
	// TODO: реализовать
}
```

После базовой реализации интервьюер продолжал ту же задачу: как развивать код для продакшена — параллельная обработка батчей, ограничение конкурентности, обработка ошибок; что произойдёт, если запустить все батчи через `errgroup` без ограничения параллелизма.
