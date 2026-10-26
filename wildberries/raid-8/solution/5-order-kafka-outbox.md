# Решение 5 — сохранение заказа и событие в Kafka

Проблема: две записи в несвязанных системах (БД и Kafka) без общей
транзакции. Между `SaveOrder` и `SendOrderCreatedEvent` процесс может упасть
или Kafka оказаться недоступна — тогда заказ в БД есть, а событие потеряно:
кто должен реагировать на заказ (оплата, склад, уведомления), о нём не
узнает. Симметричный баг при обратном порядке: событие ушло, а заказ не
сохранился — подписчик прочитает несуществующий заказ. Перезапросы и ретраи
`SendOrderCreatedEvent` не решают: ретраи дают повторы (нужен идемпотентный
потребитель), а после падения процесса некому даже начать ретрай.

Каноническое решение — transactional outbox:

1. В той же транзакции БД, что вставляет заказ, вставить строку в таблицу
   `outbox (event_id, topic, payload, created_at)` со статусом «не
   отправлено». Атомарность даёт БД: заказ без события или событие без
   заказа невозможны.
2. Отдельный процесс публикует записи outbox в Kafka и помечает их
   отправленными (или удаляет): реализация — periodic poller по таблице или
   CDC (Debezium читает WAL — минус лишняя нагрузка на БД и задержка ниже).
3. Событие помечается уникальным `event_id`/ключом, потребитель
   идемпотентен (дедупликация по ключу), потому что at-least-once допускает
   повторы доставки.

```go
func SaveOrder(ctx context.Context, tx *sql.Tx, order Order) (Order, error) {
	if err := tx.QueryRowContext(ctx,
		`INSERT INTO orders(product, quantity, status) VALUES($1,$2,$3) RETURNING id`,
		order.Product, order.Quantity, order.Status).Scan(&order.ID); err != nil {
		return order, err
	}
	payload, _ := json.Marshal(order)
	_, err := tx.ExecContext(ctx,
		`INSERT INTO outbox(event_id, topic, payload) VALUES($1,$2,$3)`,
		uuid.NewString(), "orders.created", payload)
	return order, err
}
// обе вставки завершаются одним tx.Commit(); публикацию делает poller/CDC
```

Альтернативы, которые стоит называть, но выбирать осознанно:

- «сначала Kafka, потом БД» + запись в БД из потребителя события — просто
  переносит ту же проблему в другое место;
- Kafka-transactions не помогают: транзакция Kafka не охватывает PostgreSQL;
- local-first с таблицей-статусом и компенсацией (SAGA) — тяжелее outbox,
  оправдан только когда заказ — часть большой распределённой цепочки.
