-- Задача: храним информацию о книгах (название, автор, ISBN, кому выданы,
-- когда выданы) и о пользователях (ФИО). Составить таблицы и написать запрос,
-- чтобы получить книги, которые кто-то читает, и их читателей.
--
-- Схема (по решению с собеседования):
-- author(id bigint, name varchar, birthday int)
-- author_book(id_author bigint FK -> author.id, id_book bigint FK -> book.id)
-- book(id bigint, name varchar, isbn int, reserve_date timestamptz, reader_id bigint FK -> user.id, NULL)
-- user(id bigint, name varchar)

-- Книги, которые кто-то читает, и их читатели:
SELECT u.name, b.name
FROM book b
INNER JOIN "user" u ON u.id = b.reader_id;
