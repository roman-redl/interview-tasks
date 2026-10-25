package main

import (
	"fmt"
	"sync"
	"time"
)

var cacheStore sync.Map

const timeout = time.Second

type result struct {
	data interface{}
	err  error
}

// GetData возвращает данные из getter, страхуя кешом
// в случае ошибки. Проблема: getter может отвечать очень долго.
// Задача: в случае ответа дольше timeout отдавать кеш,
// и правильно обработать кейс, когда getter всегда отвечает долго.
func GetData(key string, getter func() (interface{}, error)) (interface{}, error) {
	ch := make(chan result)

	go func() {
		data, err := getter()
		ch <- result{
			data: data,
			err:  err,
		}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			fmt.Printf("Getter result err: %s", res.err)

			if data, ok := cacheStore.Load(key); ok {
				return data, nil
			}

			return nil, res.err
		}

		cacheStore.Store(key, res.data)

		return res.data, nil
	case <-time.After(timeout):
		fmt.Printf("Getter result timeout")

		if data, ok := cacheStore.Load(key); ok {
			return data, nil
		}
	}

	return nil, nil
}

func main() {
	data, err := GetData("key", func() (interface{}, error) {
		return "some data", nil
	})

	fmt.Println(data, err)
}
