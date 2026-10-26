# Ви.тех — заход 1, техническое интервью

- Даты: 2026-03-25
- Формат: онлайн, задачи на знание Go + теория

## 1. append и запас ёмкости слайса

Что выведет программа?

```go
a := make([]int, 3, 5)
a[0], a[1], a[2] = 1, 2, 3

b := append(a, 4)
c := append(a, 5)

fmt.Println(a, b, c)
```

## 2. Рефакторинг CreateOrder

На экране дан код сервиса создания заказа (реализация OrderClient на экране отсутствовала — ниже заглушка). Требуется найти проблемы реализации и предложить рефакторинг.

```go
package main

import (
	"database/sql"
	"errors"
)

type ProductService struct {
	db          *sql.DB
	orderClient *OrderClient
}

type Product struct {
	ID       int64   // id
	Name     string  // наименование
	Price    float64 // цена
	Quantity int     // количество
}

func (s *ProductService) CreateOrder(product_id int64) error {
	var product Product

	row := s.db.QueryRow(
		"SELECT * FROM products WHERE id = ?",
		product_id,
	)

	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Quantity); err != nil {
		return err
	}

	if product.Quantity == 0 {
		return errors.New("product is out of stock")
	}

	_, err := s.db.Exec(
		"UPDATE products SET quantity = $1 WHERE id = $2",
		product.Quantity-1,
		product_id,
	)
	if err != nil {
		return err
	}

	return s.orderClient.CreateOrder(product.ID)
}

type OrderClient struct{}

func (c *OrderClient) CreateOrder(productID int64) error {
	return nil
}
```
