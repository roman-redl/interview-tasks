-- Задача 3: имя и город каждого человека из Person; если города нет — null.
-- Проверено на PostgreSQL с тестовыми данными:
--   Person: (1, Alice), (2, Bob), (3, Carol)
--   Address: (100, 1, NYC), (101, 3, LA)
--   → Alice/NYC, Bob/null, Carol/LA

SELECT p.name,
       a.city
FROM person p
         LEFT JOIN address a ON p.personid = a.personid;
