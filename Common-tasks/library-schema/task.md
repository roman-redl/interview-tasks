# Модель библиотеки

Описать структуру БД библиотеки: в библиотеке есть книги, у книг авторы, посетители могут брать книги на руки. Физически книга только одна и может быть только у одного читателя.

```sql
visitors
  id primary key
  name varchar(255)

books
  id primary key
  name varchar(255)
  visitor_id foreign key

authors
  id primary key
  name varchar(255)

books_authors
  id_book foreign key (books)
  id_author foreign key (authors)
  primary key (id_book, id_author)
```

Написать запрос: выбрать названия всех книг в библиотеке, у которых больше 3 авторов (вариант: топ самых читаемых авторов).

Встречена: Ozon (2023-09, [заход 2](../../ozon/raid-2/interview/task.md)), Wildberries Bank (2025-12-17, [задача 1](../../wildberries/raid-8/task.md)), Ozon Bank (2026-05-07 — [расширенная версия](../../ozon/raid-7/interview/task.md): книги на руках, срез на дату, топ авторов).
