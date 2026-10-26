# Ozon — заход 5, техническое интервью

- Даты: 2024-08-07 (дата по календарю)
- Формат: Go (конкурентная обработка запросов) + SQL

## 1. Обработка ссылок с остановкой по сигналу

Требуется:

1. произвести обработку ссылок — отправить GET-запросы на URL и выводить в консоль статус ответа 200 или иной;
2. распараллелить обработку ссылок;
3. реализовать остановку обработки после получения извне сигнала о прекращении работы.

```go
func main() {
	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://ёёёё",
	}
}

func getExternalSignal() chan struct{} {
	return make(chan struct{})
}
```

## 2. Порядок выполнения SQL-запроса в PostgreSQL

Определить порядок операций в ходе выполнения запроса:

```sql
SELECT
customers.name,
COUNT(order_id) as Total_orders,
SUM(order_amount) as total_spent
FROM customers
JOIN orders ON customers.id = orders.customer_id
WHERE order_date >= '2023-01-01'
GROUP BY customers.name
HAVING total_spent >= 1000
ORDER BY customers.name
LIMIT 100;
```

Ещё две задачи этапа повторяли SQL-задачи [3](../../raid-1/interview/task.md) и [4](../../raid-1/interview/task.md) захода 1 — условия повторно не вносятся.
