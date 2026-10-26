-- Задача 4: имена топ-3 читаемых авторов на данный момент.
-- «Читаемость» = сколько экземпляров книг автора сейчас на руках;
-- соавторство книги на руках даёт +1 каждому её автору.
-- Проверено на PostgreSQL: соавторы A1/A2/A3 по 2 книги, A4 — 1.

SELECT a.name
FROM authors a
         JOIN books_authors ba ON ba.author_id = a.id
         JOIN loans l ON l.book_id = ba.book_id
WHERE l.returned_at IS NULL
GROUP BY a.id, a.name
ORDER BY count(*) DESC, a.id
LIMIT 3;

-- a.id во вторичную сортировку — детерминированный порядок при равных
-- счётчиках; group by по id, а не только по имени — не склеивать
-- авторов-тёзок.
