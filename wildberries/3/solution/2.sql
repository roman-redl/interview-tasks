-- Дана таблица сотрудников компании employee
-- - id int — идентификатор сотрудника
-- - parent_id int — идентификатор его руководителя
-- - name string – имя сотрудника
--
-- Задача: вывести имена всех линейных сотрудников,
-- которые не являются руководителями.
-- Иначе говоря: вывести все листья в дереве иерархии.

SELECT name
FROM employee
WHERE id NOT IN (
    SELECT parent_id
    FROM employee
    WHERE parent_id IS NOT NULL
)