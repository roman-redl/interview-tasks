# Wildberries — заход 3, техническое интервью

- Даты: 2024-08-05
- Формат: Go-задача на многопоточность + SQL

## 1. Getter с таймаутом и кешом

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

var cacheStore sync.Map

const timeout = time.Second

// GetData возвращает данные из getter, страхуя кешом
// в случае ошибки. Проблема: getter может отвечать очень долго.
// Задача: в случае ответа дольше timeout отдавать кеш,
// и правильно обработать кейс, когда getter всегда отвечает долго.
func GetData(key string, getter func() (interface{}, error)) (interface{}, error) {
    data, err := getter()
    if err == nil {
        cacheStore.Store(key, data)
        return data, nil
    }

    fmt.Printf("Getter result err: %s", err)
    if data, ok := cacheStore.Load(key); ok {
        return data, nil
    }

    return nil, err
}
```

## 2. Листья дерева иерархии

Дана таблица сотрудников компании `employee`:

- `id int` — идентификатор сотрудника
- `parent_id int` — идентификатор его руководителя
- `name string` — имя сотрудника

Задача: вывести имена всех линейных сотрудников, которые не являются руководителями. Иначе говоря: вывести все листья в дереве иерархии.
