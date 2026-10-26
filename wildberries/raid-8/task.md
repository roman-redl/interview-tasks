# Wildberries (Bank) — заход 8, техническое интервью

- Даты: 2025-12-17
- Формат: SQL (схема, запросы, индексы) + Go (слайсы, фильтрация in-place, указатели, консистентность, семафор)

## 1. Схема библиотеки и топ авторов

Описать структуру БД библиотеки: в библиотеке есть книги, у книг авторы, посетители могут брать книги на руки. Предложенная схема:

```sql
visitors
  id primary key
  name varchar(255)

books
  id primary key
  name varchar(255)
  visitor_id foreign key

authors
  id primary key
  name varchar(255)

books_authors
  id_book foreign key (books)
  id_author foreign key (authors)
  primary key (id_book, id_author)
```

Задача: вывести топ-10 самых популярных авторов. Обсуждавшийся на собеседовании вариант:

```sql
select a.name
from authors a
inner join books_authors ba ON a.id = ba.id_author
inner join books b ON b.id = ba.id_book
group by a.name
order by count(visitor_id) DESC
limit 10
```

## 2. Индекс под запрос

Дан запрос и предложенный индекс. Оценить: подойдёт ли, в каком порядке должны идти колонки.

```sql
SELECT * FROM employees WHERE age = $1 AND salary < $2 AND sex = 'male';

create index employees (male, salary, age)
```

## 3. append и mutate слайса

Что выведется?

```go
package main

import "fmt"

func main() {
	nums := make([]int, 1, 2)
	fmt.Println(nums)

	appendSlice(nums, 1024)
	fmt.Println(nums)

	mutateSlice(nums, 1, 512)
	fmt.Println(nums)
}

func appendSlice(sl []int, val int) {
	sl = append(sl, val)
}

func mutateSlice(sl []int, idx, val int) {
	sl[idx] = val
}
```

## 4. Указатель на структуру

Что выведется?

```go
package main

import "fmt"

type Person struct {
    Name string
}

func changeName(persons *Person) {
    if persons != nil {
        persons.Name = "Alice"
    }
}

func main() {
    person := &Person{
        Name: "Bob",
    }
    fmt.Println(person.Name)
    changeName(person)
    fmt.Println(person.Name)
}
```

## 5. Сохранение заказа и событие в Kafka

Найти проблему в подходе: заказ пишется в БД, затем событие отправляется в Kafka. Что сломается при сбоях и как исправить?

```go
type Order struct {
    Product  string
    Status   string
    ID       int64
    Quantity int
}

func SaveOrder(ctx context.Context, db *sql.DB, order Order) (Order, error) {
    if err := db.QueryRowContext(ctx,
        "INSERT INTO orders(product, quantity, status) VALUES($1, $2, $3) RETURNING id",
        order.Product, order.Quantity, order.Status).Scan(&order.ID); err != nil {
        return order, err
    }

    return order, nil
}

func SendOrderCreatedEvent(order Order) error {
    // пишет в кафку
    fmt.Printf("Отправка события OrderCreated: %+v\n", order)
    return nil
}
```

## 6. Группировка слайса с семафором

Найти все проблемы в коде (на входе слайс из ~10¹⁸ элементов):

```go
func proccess(abc []string) map[string][]int {
	m := make(map[string][]int)

	mu := sync.Mutex{}
	wg := sync.WaitGroup{}
	sem := NewSemaphore(5)

	wg.Add(len(abc))
	for i, el := range abc {
		sem.Acquire()
		go func() {
			mu.Lock()
			if mel, ok := m[el]; ok {
				m[el] = append(m[el], i)
			} else {
				m[el] = make([]int, 1)
				m[el][0] = i
			}
			mu.Unlock()
			wg.Done()
			sem.Release()
		}()
	}

	wg.Wait()

	return m
}

type struct semaphore {
	chan struct{} ch
}

func NewSemaphore(thr int) semaphore {
	return semaphore{
		ch: make(chan struct{}, thr),
	}
}

func (s *semaphore) Acquire() {
	s.ch <- struct{}
}

func (s *semaphore) Release() {
	<-s.ch
}
```

## 7. Фильтрация слайса на месте

Написать функцию, не создавая новый слайс; можно обрезать слайс.

```go
func filter(a []int) []int {

}
```

Примеры из условия:

```go
[1,2,3,4,5] -> [1,3,5]
[2,2,2]     -> []
[1,2,2]     -> [1]
```
