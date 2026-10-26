package main

import (
	"context"
	"fmt"
	"io"
	"strings"
)

type Service struct {
}

func (s *Service) processDataInternal(r io.Reader) error {
	return nil
}

// Недостатки исходного ProcessData:
//  1. при срабатывании таймаута горутина навсегда зависает на отправке в небуферизованный
//     errCh — утечка горутины и памяти;
//  2. канал не закрывается после записи.
//
// Правильнее сделать так:
func (s *Service) ProcessData(timeoutCtx context.Context, r io.Reader) error {
	errCh := make(chan error, 1) // Используем буферизированный канал

	go func() {
		errCh <- s.processDataInternal(r)
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}

func main() {
	s := &Service{}
	r := strings.NewReader("data")

	// Обычное завершение — обработка успевает.
	fmt.Println(s.ProcessData(context.Background(), r))

	// Контекст уже отменён — вернётся ошибка контекста; горутина не зависнет:
	// результат отправится в буфер канала и будет собран сборщиком мусора.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(s.ProcessData(ctx, r))
}
