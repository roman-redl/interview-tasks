# Wildberries — попытка 3, техническое интервью

- Даты: ~2025-07 (по дате добавления в репозиторий, точная неизвестна)
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

Решения: базовое (`solution/1/`) и расширенное с восстановлением после серии долгих ответов (`solution/1-advanced/`).

## 2. Листья дерева иерархии

Дана таблица сотрудников компании `employee`:

- `id int` — идентификатор сотрудника
- `parent_id int` — идентификатор его руководителя
- `name string` — имя сотрудника

Задача: вывести имена всех линейных сотрудников, которые не являются руководителями. Иначе говоря: вывести все листья в дереве иерархии.

Решение: `solution/2.sql`.
