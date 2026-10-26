-- Задача 7 (MySQL): вернуть 10 пользователей, у которых нет авто с car_id = 1.
-- Таблицы users, user_cars (uniq index (user_id, car_id)); у пользователя
-- может быть неограниченное количество машин.
--
-- Ниже: генерация тестовых данных (PostgreSQL: 1 млн пользователей, у каждого
-- 0–10 случайных машин) и сравнение медленного (NOT IN) и быстрого
-- (LEFT JOIN ... IS NULL) решений через EXPLAIN ANALYZE.

CREATE TABLE users
(
    id   serial PRIMARY KEY,
    name varchar(255) NOT NULL
);

CREATE TABLE user_cars
(
    user_id integer NOT NULL REFERENCES users (id),
    car_id  integer NOT NULL,
    PRIMARY KEY (user_id, car_id)
);

BEGIN;
-- Генерация 1 миллиона пользователей
INSERT INTO users (name)
SELECT 'User_' || GENERATE_SERIES(1, 1000000);

-- Генерация случайных машин для пользователей с проверкой уникальности
-- Каждый пользователь может иметь от 0 до 10 машин
DO $$
    DECLARE
        user_id INTEGER;
        car_id INTEGER;
    BEGIN
        FOR user_id IN 1..1000000 LOOP
                FOR car_id IN 1..(1 + floor(random() * 10)) LOOP
                        BEGIN
                            INSERT INTO user_cars (user_id, car_id)
                            VALUES (user_id, (random() * 10)::integer + 1);
                        EXCEPTION
                            WHEN unique_violation THEN
                            -- Игнорируем дублирующиеся записи
                        END;
                    END LOOP;
            END LOOP;
    END $$;

COMMIT;


-- Медленное решение
EXPLAIN ANALYZE
SELECT u.id, u.name
FROM users AS u
WHERE u.id NOT IN (SELECT users.id
                   FROM users
                            INNER JOIN user_cars ON users.id = user_cars.user_id
                   WHERE user_cars.car_id = 1)
LIMIT 10;

-- Быстрое решение
EXPLAIN ANALYZE
SELECT u.id, u.name
FROM users u
         LEFT JOIN user_cars uc ON u.id = uc.user_id AND uc.car_id = 1
WHERE uc.user_id IS NULL
LIMIT 10;