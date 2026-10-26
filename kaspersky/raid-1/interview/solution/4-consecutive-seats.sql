-- Задача 4: все места, входящие в последовательность свободных длиной ≥ 2,
-- по возрастанию seat_id.
-- Проверено на PostgreSQL с тестовыми данными (1,0,1,1,1) → 3,4,5.

-- Вариант 1: оконные функции — место свободно и свободен сосед слева или справа.
SELECT seat_id
FROM (
    SELECT seat_id,
           free,
           LAG(free)  OVER (ORDER BY seat_id) AS prev_free,
           LEAD(free) OVER (ORDER BY seat_id) AS next_free
    FROM cinema
) t
WHERE free
  AND (COALESCE(prev_free, false) OR COALESCE(next_free, false))
ORDER BY seat_id;

-- Вариант 2: self-join — место свободно и есть свободное место с seat_id ± 1.
SELECT DISTINCT c1.seat_id
FROM cinema c1
         JOIN cinema c2
              ON (c1.seat_id = c2.seat_id - 1 OR c1.seat_id = c2.seat_id + 1)
                  AND c2.free
WHERE c1.free
ORDER BY c1.seat_id;
