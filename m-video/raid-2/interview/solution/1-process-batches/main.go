package main

import (
	"context"
	"fmt"
)

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
	return &Service{saver: saver, notifyer: notifier}
}

// ProcessBatches режет вход на батчи по 50 элементов и для каждого
// батча сохраняет и отправляет уведомление. Ошибки оборачиваются с
// границами батча; отмена контекста прерывает до начала очередного
// батча (Save/Notify сами обязаны уважать ctx, если ходят в сеть).
func (s *Service) ProcessBatches(ctx context.Context, items []string) error {
	const batchSize = 50

	for start := 0; start < len(items); start += batchSize {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("batches: %w", err)
		}
		end := start + batchSize
		if end > len(items) {
			end = len(items)
		}
		batch := items[start:end]

		if err := s.saver.Save(batch); err != nil {
			return fmt.Errorf("save batch [%d:%d]: %w", start, end, err)
		}
		if err := s.notifyer.Notify(batch); err != nil {
			return fmt.Errorf("notify batch [%d:%d]: %w", start, end, err)
		}
	}
	return nil
}

// Демо: 120 элементов -> батчи 50 + 50 + 20.
type demoSaver struct{ calls int }

func (d *demoSaver) Save(data []string) error {
	d.calls++
	fmt.Printf("save #%d: %d items\n", d.calls, len(data))
	return nil
}

type demoNotifier struct{ calls int }

func (d *demoNotifier) Notify(data []string) error {
	d.calls++
	fmt.Printf("notify #%d: %d items\n", d.calls, len(data))
	return nil
}

func main() {
	sv, nt := &demoSaver{}, &demoNotifier{}
	svc := New(sv, nt)

	items := make([]string, 120)
	for i := range items {
		items[i] = fmt.Sprintf("item-%d", i)
	}

	if err := svc.ProcessBatches(context.Background(), items); err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("saves:", sv.calls, "notifies:", nt.calls) // 3 и 3
}
