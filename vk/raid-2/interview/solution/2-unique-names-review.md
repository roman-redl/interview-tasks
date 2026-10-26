# Решение 2 — ревью: уникальные имена пользователей

Проблемы по возрастанию:

1. **Паника: запись в nil map.** `userFirstNamesMap` объявлена, но не
   инициализирована — `userFirstNamesMap[userFirstName] = true` паникует
   `assignment to entry in nil map`. Фикс: `make(map[string]bool)` (или
   `make(map[string]struct{}, len(userFirstNames))`).
2. **Цикл без закрывающей скобки** — код вообще не скомпилируется
   (фрагмент; но даже по фрагменту: тело цикла не закрыто, `return`
   отсутствует).
3. **Объём в память**: тащить ВСЕ имена пользователей проекта строками,
   чтобы потом дедупнуть в Go — при больших проектах это мегабайты и
   время. Дедуп — задача SQL: `SELECT DISTINCT first_name FROM ...`;
   ещё лучше — не `SELECT *`-стайл: явно перечислить колонку, добавить
   индекс под фильтр проекта.
4. **Работа с Rows**: в интерфейсе есть `Close()`, но по фрагменту не
   видно вызова (обычная ошибка джуна): не закрытые rows держат
   соединение пула; `Next()`-цикл обязан проверять `rows.Err()` после
   завершения — ошибка чтения в середине молча даёт неполный результат.
5. **Глобальная `var conn Conn`**: зависимости через глобальную
   переменную — не тестируется и не потокобезопасно инициализируется;
   передавать Conn параметром/конструктором хранилища.
6. **Кодстайл**: смешение возвратов `nil, err`; имя
   `GetUniqueUserFirstNamesByProjectName` — слишком длинно, а два слоя — «получить имена» и «дедупнуть» — правильнее слить в один
   запрос DISTINCT; `if _, ok := map[...]` — идиоматично, но с
   `map[string]bool` лучше просто читать значение.

Исправленный вид:

```go
func GetUniqueUserFirstNamesByProject(ctx context.Context, projectName string) ([]string, error) {
	rows, err := conn.Query(ctx,
		`SELECT DISTINCT first_name FROM users u
		 JOIN projects p ON p.id = u.project_id
		 WHERE p.name = $1`, projectName)
	if err != nil {
		return nil, fmt.Errorf("query user names: %w", err)
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			return nil, fmt.Errorf("scan first_name: %w", err)
		}
		names = append(names, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}
	return names, nil
}
```

Итог ревью: паника nil-map и незакрытый цикл — блокеры; остальное —
«дедуп делает БД, не горутина», гигиена rows и отказ от глобального
состояния.
