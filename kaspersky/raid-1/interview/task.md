# Kaspersky — заход 1, техническое интервью

- Даты: 2026-03-26
- Формат: онлайн, разбор опыта и серия мини-задач «что выведется» на Go (слайсы, каналы, интерфейсы, структуры) плюс две SQL-задачи; код показывался на экране, ответ голосом

## 1. Слайсы над общим массивом

Дан код. Что он выведет?

```go
func main() {
	original := []int{1, 2, 3, 4, 5}
	s1 := original[1:4]
	s2 := original[2:]
	s1[0] = 100
	fmt.Println(s1)
	fmt.Println(s2)
	fmt.Println(original)
}
```

## 2. Append и capacity: два слайса

Что выведет программа?

```go
func main() {
	s1 := make([]int, 0, 3)
	s1 = append(s1, []int{1, 2, 3}...)
	s2 := append(s1, []int{4, 5, 6}...)
	s2[0] = 10
	if s1[0] != s2[0] {
		fmt.Println("elements are not euqal")
		return
	}
	fmt.Println("elements are equal")
}
```

Следом в том же коде меняли ёмкость при создании слайса: `make([]int, 0, 6)`. Что выведется теперь? И как изменить исходный код, чтобы вывелось `elements are equal`?

## 3. SQL: люди без адресов

Даны таблицы `Person` и `Address`:

```text
Таблица: Person
+---------------+----------+
| Column Name   | Type     |
+---------------+----------+
| personId      | int      |
| Name          | varchar  |
+---------------+----------+
personId - primary key


Таблица: Address
+---------------+----------+
| Column Name   | Type     |
+---------------+----------+
| addressId     | int      |
| personId      | int      |
| city          | varchar  |
+---------------+----------+
addressId - primary key.
```

Напишите запрос, выводящий имя и город каждого человека из таблицы `Person`. Если города для `personId` не находится, должен быть `null`.

## 4. SQL: соседние свободные места

Таблица `Cinema`:

```text
Таблица: Cinema
+-------------+------+
| Column Name | Type |
+-------------+------+
| seat_id     | int  |
| free        | bool |
+-------------+------+
seat_id это auto-increment
```

Каждая строка таблицы указывает, свободно ли i-е место или нет: `1` означает свободное место, а `0` — занятое. Найти все последовательно доступные места в кинотеатре (минимальная длина последовательности — 2 места). Вернуть результирующую таблицу, отсортированную по `seat_id` в порядке возрастания.

Пример:

```text
Input:
+---------+------+
| seat_id | free |
+---------+------+
| 1       | 1    |
| 2       | 0    |
| 3       | 1    |
| 4       | 1    |
| 5       | 1    |
+---------+------+
Output:
+---------+
| seat_id |
+---------+
| 3       |
| 4       |
| 5       |
+---------+
```

## 5. Ресивер по значению

Что выведет программа и что нужно изменить, чтобы метод действительно увеличивал переменную из `main`?

```go
type Count int

func (c Count) Increment() {
	c++
}

func main() {
	var count Count
	count.Increment()
	fmt.Print(count)
}
```

## 6. Select и закрытые каналы

Что произойдёт при запуске программы?

```go
func main() {
	msgChan := make(chan string, 1)
	close(msgChan)
	stopChan := make(chan struct{})
	close(stopChan)

	select {
	case msgChan <- "msg":
		fmt.Println("msg sent")
	case <-stopChan:
		fmt.Println("stop signal received")
	}
}
```

## 7. Select, буферный канал и nil-канал

Что выведется?

```go
func main() {
	c := make(chan int, 1)
	for done := false; !done; {
		select {
		default:
			done = true
			fmt.Println(1)
		case <-c:
			fmt.Println(2)
			c = nil
		case c <- 1:
			fmt.Println(3)
		}
	}
}
```

## 8. Nil-указатель внутри интерфейса

Что выведется? Как проверить, что внутри интерфейса не лежит значение?

```go
func main() {
	var s *string
	fmt.Println(s == nil) // (1)

	var i interface{}
	fmt.Println(i == nil) // (2)

	i = s
	fmt.Println(i == nil) // (3)
}
```

## 9. Копии значения: интерфейс и горутина

Что выведется?

```go
type Speaker interface {
	Speak() string
}

type Person struct {
	name  string
	alias string
}

func (p Person) Speak() string {
	return "Hi, I'm " + p.name
}

func main() {
	var s Speaker
	p := Person{name: "Alice"}
	s = p

	p.name = "Bob"

	fmt.Println(s.Speak()) // (1)

	func() {
		s := p
		fmt.Println(s.Speak()) // (2)
	}()
	go func(p Person) {
		time.Sleep(10 * time.Millisecond)
		fmt.Println(p.Speak()) // (3)
	}(p)

	p.name = "Charlie"

	time.Sleep(20 * time.Millisecond)
}
```

## 10. Размер структур и выравнивание

Что выведется и почему?

```go
type Car1 struct {
	numOfWheels int
	newOnMarket bool
	color       int
}

type Car2 struct {
	numOfWheels int
	newOnMarket bool
	fromChina   bool
	fromEurope  bool
	fromUSA     bool
	color       int
}


func main() {
	var c1 Car1
	var c2 Car2
	if unsafe.Sizeof(c1) == unsafe.Sizeof(c2) {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}
}
```
