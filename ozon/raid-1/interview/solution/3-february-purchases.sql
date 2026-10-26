CREATE TABLE userr
(
    id     serial PRIMARY KEY,
    first_name  varchar(255),
    last_name varchar(255)
);

INSERT INTO userr (first_name, last_name)
VALUES ('f1', 's1'),
       ('s2', 's2');

CREATE TABLE purchase
(
    sku     varchar(255),
    price   integer,
    user_id integer REFERENCES userr (id),
    date    timestamp
);

INSERT INTO purchase (sku, price, user_id, date)
VALUES ('11111', 5500, 1, '2021-02-15'),
       ('22222', 4000, 1, '2021-02-14'),
       ('33333', 8000, 2, '2021-03-01'),
       ('44444', 400, 2, '2021-03-02');

CREATE TABLE ban_list
(
    user_id integer REFERENCES userr (id),
    date    timestamp
);

INSERT INTO ban_list (user_id, date)
VALUES (1, '2021-03-08');

-- Задача 3: список уникальных клиентов и купленных ими товаров в феврале

SELECT DISTINCT userr.first_name, userr.last_name, purchase.sku
FROM userr
         INNER JOIN purchase ON userr.id = purchase.user_id
WHERE purchase.date >= '2021-02-01' AND purchase.date < '2021-03-01';
