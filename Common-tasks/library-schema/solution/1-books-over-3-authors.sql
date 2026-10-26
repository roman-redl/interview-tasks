-- Задача 1: названия всех книг, у которых больше 3 авторов.
-- Проверено на PostgreSQL с тестовыми данными:
--   BookA (4 автора, на руках), BookB (2), BookC (1, на полке), BookD (2)
--   → только BookA.

SELECT b.name
FROM books b
         JOIN books_authors ba ON b.id = ba.id_book
GROUP BY b.id, b.name
HAVING count(*) > 3;
