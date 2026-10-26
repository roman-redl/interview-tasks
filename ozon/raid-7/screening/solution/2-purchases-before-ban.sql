-- Задача 2: уникальные комбинации пользователя и SKU для покупок,
-- совершённых ДО бана этого пользователя. Не забаненные пользователи
-- в выборку не попадают (JOIN по ban_list — фильтр).
-- Проверено на PostgreSQL с данными из условия:
--   Ivan|1, Ivan|2 (у Ивана две покупки SKU 1 — DISTINCT схлопывает).

SELECT DISTINCT u.firstname,
       p.sku
FROM purchase p
         JOIN "user" u ON u.id = p.user_id
         JOIN ban_list b ON b.user_id = p.user_id
WHERE p.date < b.date_from
ORDER BY u.firstname, p.sku;

-- Примечания:
-- - если по условию выводить id пользователя, а не имя — заменить первую
--   колонку на u.id (в данных есть тёзки Anna Petrova, DISTINCT по имени
--   их склеил бы при выводе имени);
-- - строгое < : покупка в день бана бана «до бана» не пережила.
