-- Задача 3: пользователи с суммой покупок больше 5000 р.
-- Проверено на PostgreSQL с данными из условия:
--   1|Ivan|Petrov|15200, 2|Anna|Petrova|8400.

SELECT u.id,
       u.firstname,
       u.lastname,
       sum(p.price) AS total
FROM "user" u
         JOIN purchase p ON p.user_id = u.id
GROUP BY u.id, u.firstname, u.lastname
HAVING sum(p.price) > 5000
ORDER BY u.id;
