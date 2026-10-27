# Rambler & Co — заход 1, техническое интервью

- Даты: 2025-07-23
- Формат: онлайн, архитектурные мини-задачи + live-coding на Go

## 1. Внешняя сортировка страниц

Предположим, у нас есть 4 страницы данных, и мы не можем одновременно загрузить их в оперативную память. Каждая страница содержит символы, которые нужно записать в итоговый один файл, но в отсортированном виде. Каким образом можно решить эту задачу?

## 2. Чтение из БД с кешем

Есть функция получения данных из «базы». Необходимо придумать вариант, как правильно читать данные из БД: данные запрашиваются часто и много. Для примера используется цикл, который пытается получить данные:

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID   int
	Name string
}

var (
	dbHits int
)

func getUserInfoFromDB(id int) User {
	time.Sleep(100 * time.Millisecond)
	dbHits++
	return User{ID: id, Name: fmt.Sprintf("User-%d", id)}
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			user, ok := getUserInfoFromCache(i)
		}(42)
	}
	wg.Wait()
	fmt.Println("DB hits:", dbHits)
}
```
