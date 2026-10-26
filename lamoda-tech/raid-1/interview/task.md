# Lamoda Tech — заход 1, техническое интервью

- Даты: 2026-05-05
- Формат: онлайн, live-coding (Go) + SQL + разбор кода

## 1. Произведение всех остальных элементов

Условие — комментарий в редакторе; требуется реализовать функцию.

```go
// Есть слайс целых чисел, нужно написать функцию которая вернет слайс такого же размера,
// где на каждой i-ой позиции будет произведение всех элементов кроме i-го
// [1, 2, 3] => [2*3, 1*3, 1*2] => [6, 3, 2]

func MultOther(in []int) []int {

}
```

Уточнения по ходу: как решение ведёт себя при нулях в слайсе; корнер-кейсы по размеру слайса (0, 1, 2 элемента).

## 2. Параллельные HTTP-запросы с ограничением конкурентности

Условие — код в редакторе: список URL для обхода и пустая функция, которую требуется реализовать. Строка из шести точек в списке URL — так на экране, реальный список длиннее показанного.

```go
// Требуется обойти все url и получить слайс с кодами ответов в том же порядке
// Хотим делать параллельно но не больше k запросов одновременно
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

var urls = []string{
	"https://www.lamoda.ru/p/mp002xw0lvkd/clothes-tomollyfromjames-plate/",
	"https://www.lamoda.ru/p/mp002xw14uf2/clothes-tomollyfromjames-plate/",
	"https://www.lamoda.ru/p/rtladr746901/clothes-iceberg-plate/",
	"https://www.lamoda.ru/p/mp002xw18h9d/clothes-victoriaveisbrut-plate/",
	"https://www.lamoda.ru/p/mp002xw004x4/clothes-clanvi-plate/",
	"https://www.lamoda.ru/p/mp002xw0zfxy/clothes-glvr-plate/",
	"https://www.lamoda.ru/p/mp002xw0slmg/clothes-snezhnayakoroleva-plate-kozhanoe/",
	"https://www.lamoda.ru/p/mp002xw132c3/clothes-auranna-plate/",
	......
}

func crawl(urls []string, k int) []int {

}

func main() {
	result := crawl(urls, 5)
	fmt.Println("All done")
}
```

Уточнения по ходу: входной слайс может быть очень большим — решение с горутиной на каждый URL и планировщиком/GC не устраивает; канал задач с буфером длиной всего слайса избыточен по памяти; какой размер буфера выбрать, если небуферизированный канал снижает производительность, а буфер длиной `len(urls)` — слишком дорогой.

## 3. SQL: клиенты и их корзины

Условие и структура таблиц — код в редакторе; требуется написать запрос.

```sql
CREATE TABLE customer
(
    id      INTEGER PRIMARY KEY,
    email   VARCHAR(100)    NOT NULL,
    country CHAR(2)         NOT NULL
);

CREATE TABLE cart_item
(
    id           INTEGER PRIMARY KEY,
    customer_id  INTEGER         NOT NULL,
    title        VARCHAR(20)     NOT NULL,
    amount       INTEGER         NOT NULL,
    price        INTEGER         NOT NULL
);


-- Вывести построчно всех кастомеров(id, email)
-- и все элементы их корзины (title, amount)
```

## 4. SQL: топ клиентов по сумме корзины

На тех же таблицах `customer` / `cart_item`. Условие — комментарий в редакторе:

```sql
-- Вывести топ-10 покупателей(id, email, total_sum)
-- по общей стоимости товаров в корзине
```

Общая стоимость корзины — сумма `amount * price` по всем позициям клиента. Добавленные в обсуждении условия: оставить только клиентов из России и только тех, у кого сумма корзины больше 1000 условных рублей.

## 5. Составной индекс: какие запросы его используют

Код в редакторе: создание индекса и два запроса. Для каждого SELECT определить, будет ли использоваться индекс, и объяснить почему. Оформление — как на экране: вопрос задан комментарием через `//`, во втором запросе нет точки с запятой.

```sql
// Будет ли использоваться индекс
CREATE INDEX myIdx ON carts(sku, country, customer_id);
SELECT * FROM carts WHERE sku = 192 AND country = 'ru';
SELECT * FROM carts WHERE country = 'ru' AND customer_id = 10
```

## 6. Хеш-индекс и диапазонный запрос

Код в редакторе; ответить на вопрос из комментария — будет ли использоваться хеш-индекс в запросе с диапазонным условием.

```sql
-- Будет ли использоваться хеш индекс?
CREATE TABLE users(username text, password text, age integer);
CREATE INDEX h_idx ON users USING hash(age);
SELECT * FROM users WHERE age > 18;
```

Уточнение по ходу: как проверить фактическое использование индекса.

## 7. Ревью кода: списание баланса со счёта

Дан код (DDL таблицы и функция списания). Требуется провести ревью с точки зрения работы с базой данных: найти проблемы и предложить исправления. Стиль кода (нейминг, деферры) не интересует — только работа с БД.

```go
-- Юзеры и их балансы
CREATE TABLE users
(
    user_id INTEGER PRIMARY KEY,
    balance REAL    NOT NULL
);
// Снятие денег с баланса пользователя
func withdrawBalance(userId int32, amount float32) float32 {
	tx := BeginTransaction("ISOLATION LEVEL READ COMMITTED")

	oldBalance := tx.exec("SELECT balance FROM users WHERE user_id = $1", userId)
	if oldBalance >= amount  {
		tx.exec("UPDATE users SET balance = balance - $2 WHERE user_id = $1", userId, amount)
	}
	newBalance := tx.exec("SELECT balance FROM users WHERE user_id = $1", userId)

	tx.Commit()

	return newBalance
}
```

Обсуждались варианты исправления гонки списания: атомарный `UPDATE ... WHERE balance >= $1`, `SELECT FOR UPDATE`, поднятие уровня изоляции, блокировка на уровне приложения.
