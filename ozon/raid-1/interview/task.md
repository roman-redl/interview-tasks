# Ozon — заход 1, техническое интервью

- Даты: 2023-07-24 (дата по календарю)
- Формат: Go-задачи «что выведется» + SQL

## 1. Что выведется: append к слайсу

Что выведет программа?

```go
func a() {
    x := []int{}
    x = append(x, 0)
    x = append(x, 1)
    x = append(x, 2)
    y := append(x, 3)
    z := append(x, 4)
    fmt.Println(y, z)
}

func main() {
    a()
}
```

## 2. Что выведется: чтение из двух каналов

Что выведет программа?

```go
func main() {
    timeStart := time.Now()
    _, _ = <-worker(), <-worker()
    println(int(time.Since(timeStart).Seconds()))
}

func worker() chan int {
    ch := make(chan int)
    go func() {
        time.Sleep(3 * time.Second)
        ch <- 1
    }()
    return ch
}
```

## 3. SQL: покупки в феврале

Даны таблицы:

```text
// user
id | firstname | lastname | birth
1  | Ivan      | Petrov   | 1996-05-01
2  | Anna      | Petrova  | 1999-06-01

// purchase
sku   | price | user_id | date
11111 | 5500  | 1       | 2021-02-15
22222 | 4000  | 1       | 2021-02-14
33333 | 8000  | 2       | 2021-03-01
44444 | 400   | 2       | 2021-03-02

// ban_list
user_id | date
1       | 2021-03-08
```

Вывести список уникальных клиентов и купленных ими товаров в феврале.

## 4. SQL: покупка дороже 5000 без бана

На тех же таблицах вывести список уникальных клиентов, совершивших покупку товаров на сумму больше 5000 и не имеющих бана.
