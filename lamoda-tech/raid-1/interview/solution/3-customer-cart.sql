-- Задача 3: построчно все клиенты (id, email) и все элементы их корзин
-- (title, amount). LEFT JOIN: клиент без товаров в корзине всё равно
-- попадает в выборку — с NULL вместо товара.
-- Проверено на PostgreSQL: клиент без корзины выводится строкой с NULL.

SELECT c.id,
       c.email,
       ci.title,
       ci.amount
FROM customer c
         LEFT JOIN cart_item ci ON ci.customer_id = c.id
ORDER BY c.id, ci.id;
