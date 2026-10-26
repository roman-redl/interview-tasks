-- Задача 4: топ-10 покупателей (id, email, total_sum) по общей стоимости
-- товаров в корзине; только Россия (country = 'RU') и только сумма
-- корзины строго больше 1000.
-- Проверено на PostgreSQL с краевыми данными: клиент с суммой ровно 1000
-- и клиент с суммой 800 отфильтрованы, non-RU не попали, пустые корзины
-- в топ не попадают (JOIN).

SELECT c.id,
       c.email,
       sum(ci.amount * ci.price) AS total_sum
FROM customer c
         JOIN cart_item ci ON ci.customer_id = c.id
WHERE c.country = 'RU'
GROUP BY c.id, c.email
HAVING sum(ci.amount * ci.price) > 1000
ORDER BY total_sum DESC
LIMIT 10;

-- Замечания:
-- - country фильтруется в WHERE (до группировки), условие по агрегату
--   sum(...) — только в HAVING: это его единственное место;
-- - пустая корзина: INNER JOIN исключает клиента с NULL-суммой —
--   для «больше 1000» это и так верно, поэтому LEFT JOIN не нужен.
