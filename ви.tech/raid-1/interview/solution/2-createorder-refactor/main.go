package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrOutOfStock = errors.New("product is out of stock or does not exist")

type ProductService struct {
	db          *sql.DB
	orderClient *OrderClient
}

type Product struct {
	ID       int64
	Name     string
	Price    float64
	Quantity int
}

/*
Список всех проблем в коде:
1. Между Select и Update гонка, нет транзакции
	Если два пользователя одновременно попытаются купить последний товар (Quantity = 1),
	оба потока прочитают Quantity = 1, оба пройдут проверку product.Quantity == 0, и оба
	выполнят UPDATE. В итоге товар будет продан дважды, а в базе может оказаться отрицательное
	значение (если тип колонки позволяет) или просто 0, но по факту будет создан лишний заказ.
	Решение: добавить транзакцию и коммит + роллбэк,
		ещё лучше сделать Update + Returning в 1 запросе.
2. SELECT * в строке 17.
	Использовать звездочку в связке с row.Scan — плохая практика. Если в таблицу products
	добавят новую колонку или изменят их порядок, Scan упадет с ошибкой или запишет данные
	не в те поля структуры. Колонки нужно перечислять явно.
3. Хардкод SQL запросов в слое бизнес логики.
	SQL-запросы прописаны прямо в методе сервиса. Это затрудняет тестирование (нужен мок БД)
	 и поддержку. Лучше вынести работу с данными в слой Repository.
4. Операции CreateOrder и обновления количества товара должны быть в одной транзакции
	Если операция обновления количества товара провалится, то заказ не будет создан. Нужны
	компенсирующие действия, например использовать транзакции для обеспечения атомарности операций.
6. Небезопасная работа с типами (Money).
	Поле Price имеет тип float64. Для финансовых операций это табу из-за ошибок округления (например,
	0.1 + 0.2 != 0.3). Нужно использовать int64 (хранить в копейках/центах) или специальный тип decimal.
4. Нет проверки валидность product_id, на наличие товара (проверка sql.ErrNoRows), пользователя, заказа и т.д.
5. Нет логирования, метрик (latency, успешные заказы).
6. Нет тестирования.
7. Нейминг и обработка ошибок
	- Аргумент product_id нарушает конвенцию Go (camelCase). Должно быть productID.
	- Ошибки возвращаются "как есть". Лучше оборачивать их через fmt.Errorf("...: %w", err),
	чтобы сохранить контекст для логов. Ошибка errors.New("product is out of stock") не
	типизирована. В вызывающем коде будет сложно понять, почему именно упал метод, не парся строку.
	- Ошибку "out of stock" лучше вынести в глобальную переменную (Sentinel Error), чтобы
	вызывающий код мог проверить её через errors.Is.
9. Нарушение инкапсуляции: Структура Product используется и для БД, и для бизнес-логики, и,
	судя по всему, передается во внешний сервис. При изменении схемы таблицы в БД может "поплыть"
	API внешнего клиента.
10. Контекст. Не передаётся контекст, нужно ограничивать работу методов для graceful shutdown,
	а также иметь возможность ограничить долгий запрос в БД по таймауту, если запросы будут долгие
	из-за внешнего сервиса.
11. Нет идемпотентности. клиент отправил запрос, заказ создался, ответ не дошёл (timeout),
	клиент повторяет запрос👉Получишь два заказа. Решение: idempotency key.
12. Отсутствие ретрай логики. Сейчас просто падение при deadlock или сбое.
13. Нет защиты от падения order сервиса. Нужен timeout / circuit breaker.
14. Lost update. В коде делается SET quantity = $1 и передается product.Quantity - 1.
	Это классическая ошибка. Даже внутри транзакции правильнее делать
	SET quantity = quantity - 1, чтобы база сама инкрементировала значение атомарно.
*/

func (s *ProductService) CreateOrder(productID int64) error {
	var product Product

	// Атомарно уменьшаем количество, только если оно больше 0.
	// RETURNING вернет обновленные данные, если UPDATE прошел успешно.
	query := `
		UPDATE products 
		SET quantity = quantity - 1 
		WHERE id = $1 AND quantity > 0 
		RETURNING id, name, price, quantity
	`

	err := s.db.QueryRow(query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Quantity,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Если ни одна строка не обновилась, значит либо товара нет, либо quantity = 0
			return ErrOutOfStock
		}

		return fmt.Errorf("failed to update product quantity: %w", err)
	}

	return s.orderClient.CreateOrder(product.ID)
}

// 2 вариант Если логика сложнее, или база данных не поддерживает RETURNING,
// необходимо использовать транзакцию и блокировку строки на чтение/запись (SELECT ... FOR UPDATE).
func (s *ProductService) CreateOrder2(ctx context.Context, productID int64) error {
	// 1. Начинаем транзакцию
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	// Гарантируем откат в случае паники или ошибки
	defer tx.Rollback()

	var product Product
	// 2. Блокируем строку для других транзакций до завершения текущей (FOR UPDATE)
	query := `SELECT id, name, price, quantity FROM products WHERE id = $1 FOR UPDATE`

	err = tx.QueryRowContext(ctx, query, productID).Scan(
		&product.ID, &product.Name, &product.Price, &product.Quantity,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("product not found")
		}
		return fmt.Errorf("failed to select product: %w", err)
	}

	// 3. Проверяем бизнес-логику
	if product.Quantity <= 0 {
		return ErrOutOfStock
	}

	// 4. Обновляем данные
	_, err = tx.ExecContext(ctx, `UPDATE products SET quantity = quantity - 1 WHERE id = $1`, productID)
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	// 5. Коммитим изменения
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

type OrderClient struct{}

func (c *OrderClient) CreateOrder(productID int64) error {
	return nil
}
