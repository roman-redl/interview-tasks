-- Задача 2: топ-10 самых популярных авторов. Популярность = сколько
-- экземпляров книг автора сейчас на руках (visitor_id IS NOT NULL);
-- подсчёт по связке автор–книга: соавторство книги на руках даёт +1
-- каждому её автору.
-- Проверено на PostgreSQL с тестовыми данными:
--   A1 = 3 (BookA, BookB, BookD), A2 = 2, остальные по 1.

SELECT a.id,
       a.name,
       count(*) AS books_on_hands
FROM authors a
         JOIN books_authors ba ON a.id = ba.id_author
         JOIN books b ON b.id = ba.id_book
WHERE b.visitor_id IS NOT NULL
GROUP BY a.id, a.name
ORDER BY count(*) DESC, a.id
LIMIT 10;

-- На собеседовании (Wildberries Bank) обсуждался почти такой же запрос,
-- но с тремя слабыми местами:
--
--   select a.name
--   from authors a
--   inner join books_authors ba ON a.id = ba.id_author
--   inner join books b ON b.id = ba.id_book
--   group by a.name
--   order by count(visitor_id) DESC
--   limit 10
--
-- 1. group by a.name вместо a.id, a.name — два разных автора-тёзки
--    склеятся в одну строку.
-- 2. count(visitor_id) без WHERE — вместо фильтра считается на лету:
--    работает (count по столбцу игнорирует NULL), но менее явно; проще
--    фильтровать в WHERE, тогда count(*) считает только книги на руках.
-- 3. В списке выбора нет самого счётчика и нет вторичного ключа
--    сортировки — порядок авторов с равным числом книг недетерминирован.
