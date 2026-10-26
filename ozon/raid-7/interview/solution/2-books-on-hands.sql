-- Задача 2: названия всех книг, которые на руках.
-- Проверено на PostgreSQL: с тестовыми данными (см. 1-library-schema.sql)
-- возвращает B1, B2, B3.

SELECT DISTINCT b.name
FROM books b
         JOIN loans l ON l.book_id = b.id
WHERE l.returned_at IS NULL
ORDER BY b.name;

-- DISTINCT страхует от нескольких открытых строк loans у одной книги
-- (грязные данные при отсутствии ограничения из решения 1).
-- Вариант с денормализованным books.visitor_id:
--   SELECT name FROM books WHERE visitor_id IS NOT NULL;
