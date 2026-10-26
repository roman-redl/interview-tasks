-- Задача 1: схема библиотеки. Три сущности — Автор, Книга, Читатель;
-- физический экземпляр книги один и может быть только у одного читателя;
-- у книги несколько авторов (M:N). Для исторических срезов (задача 5)
-- добавлена таблица выдач loans — история вместо перезаписи текущего
-- держателя.
-- DDL и все запросы ниже проверены на PostgreSQL с тестовыми данными.

CREATE TABLE visitors
(
    id   bigserial PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE authors
(
    id   bigserial PRIMARY KEY,
    name text NOT NULL
);

CREATE TABLE books
(
    id   bigserial PRIMARY KEY,
    name text NOT NULL
    -- вариант денормализации: visitor_id bigint REFERENCES visitors —
    -- «кто держит сейчас», без JOIN для запроса 2. Плата — согласованность
    -- с loans при каждой выдаче/возврате. Чистый вариант — выводить
    -- держателя из loans (см. ниже) и visitor_id не хранить.
);

CREATE TABLE books_authors
(
    book_id   bigint NOT NULL REFERENCES books,
    author_id bigint NOT NULL REFERENCES authors,
    PRIMARY KEY (book_id, author_id)
);

CREATE TABLE loans
(
    id          bigserial PRIMARY KEY,
    book_id     bigint NOT NULL REFERENCES books,
    visitor_id  bigint NOT NULL REFERENCES visitors,
    taken_at    date NOT NULL, -- дата выдачи
    returned_at date           -- NULL = книга на руках сейчас
);

-- Почему история нужна: без loans «книга у читателя» — это перезаписываемое
-- поле, и на любой вопрос «кто держал книгу в прошлом / сколько было на
-- руках 1 мая» ответить нечем. Каждая выдача — отдельная строка loans:
-- повторные выдачи не затирают друг друга. Ограничение «у книги не больше
-- одной открытой выдачи» (частичный уникальный индекс) — по вкусу stricter:
-- CREATE UNIQUE INDEX ON loans (book_id) WHERE returned_at IS NULL;
