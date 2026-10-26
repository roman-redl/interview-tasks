# Точка Банк — заход 1, техническое интервью

- Даты: 2026-02-24
- Формат: онлайн, задачи на знание Go + SQL

## 1. Заменить точку на «ы»

```go
package main

func main() {
    s := "Привет."
    // s == "Приветы" 
}
```

## 2. Что выведется: присваивание nil интерфейсу

Написать, что выведется.

```go
package main

import "fmt"

type MyStruct struct {
    Value string
}

type MyInterface interface {
    GetValue() string
}

func (m *MyStruct) GetValue() string {
    return m.Value
}

func ToNil(obj MyInterface) {
    if obj == nil {
        fmt.Printf("obj = nil")
        return
    }
    obj = nil
}

func main() {
    myStruct := &MyStruct{Value: "test"}

    ToNil(myStruct)
    fmt.Printf("myStruct: %+v \n", myStruct)

    myStruct = nil

    ToNil(myStruct)
    fmt.Printf("myStruct: %+v \n", myStruct)
}
```

## 3. Оптимальный и покрывающий индекс

Написать, какой оптимальный индекс и какой покрывающий.

```sql
SELECT * FROM employee WHERE sex = 'm' AND salary > 300000 AND age = 20 ORDER BY created_at
```

## 4. Схема библиотеки и запрос читателей

Храним информацию о книгах (название, автор, ISBN, кому выданы, когда выданы) и о пользователях (ФИО). Составить таблицы и написать запрос, чтобы получить книги, которые кто-то читает, и их читателей.
