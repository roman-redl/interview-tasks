-- Задача 3: названия всех книг, у которых больше 3 авторов.
-- Проверено на PostgreSQL: книга с 4 авторами возвращается, с 2 — нет.
-- Базовая постановка задачи (без истории выдач) —
-- Common-tasks/library-schema, решение 1 там же.

SELECT b.name
FROM books b
         JOIN books_authors ba ON ba.book_id = b.id
GROUP BY b.id, b.name
HAVING count(*) > 3;
