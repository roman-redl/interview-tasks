-- Задача 3: имена, количество заказов, сумма — только у пользователей
-- с более чем 3 заказами, по убыванию суммы.
-- Проверено на PostgreSQL с тестовыми данными:
--   Alice 4/425.00, Carol 4/300.00 — проходят; Bob 3 — отфильтрован.

SELECT u.name,
       count(*)      AS orders_count,
       sum(o.amount) AS total_amount
FROM users u
         JOIN orders o ON o.user_id = u.id
GROUP BY u.id, u.name
HAVING count(*) > 3
ORDER BY total_amount DESC;

-- JOIN против LEFT JOIN: фильтр «больше 3 заказов» всё равно исключает
-- клиентов без заказов, так что внешнее соединение не нужно.
-- Группировка по u.id (PK) — не склеивать тёзок; name в GROUP BY для
-- строгого режима only_full_group_by.
