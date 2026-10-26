-- Задача 2 (схема и данные — в task.md): серии запросов по городам и пользователям.

-- 1. Вывести города по пользователям.
SELECT c.name AS city, u.name AS "user"
FROM cities c
         JOIN users u ON u.city_id = c.id;

-- 2. Вывести города с пользователями и без.
SELECT cities.name, users.name
FROM cities
         LEFT JOIN users ON users.city_id = cities.id;

-- 3. Вывести города с пользователями и без с сортировкой.
SELECT cities.name, users.name
FROM cities
         LEFT JOIN users ON users.city_id = cities.id
ORDER BY cities.name, users.name;

-- 4. Вывести города с пользователями и без и количеством пользователей.
SELECT cities.name, COUNT(users.id) AS users_count
FROM cities
         LEFT JOIN users ON users.city_id = cities.id
GROUP BY cities.name;

-- 5. Вывести города с пользователями и без, заменить 0 пользователей на NULL
--    и сортировать по имени.
SELECT cities.name, NULLIF(COUNT(users.id), 0) AS users_count
FROM cities
         LEFT JOIN users ON users.city_id = cities.id
GROUP BY cities.name
ORDER BY cities.name;

-- 6. То же с номером строки.
-- В исходном решении с собеседования вместо оконной функции использовалась
-- переменная @row_number() OVER (...) — так нумерация не работает; корректный
-- вариант — ROW_NUMBER():
SELECT cities.name,
       NULLIF(COUNT(users.id), 0)     AS users_count,
       ROW_NUMBER() OVER (ORDER BY cities.name) AS row_num
FROM cities
         LEFT JOIN users ON users.city_id = cities.id
GROUP BY cities.name
ORDER BY cities.name;
