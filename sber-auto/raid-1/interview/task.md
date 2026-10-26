# СберАвто — заход 1, техническое интервью

- Даты: неизвестна (~2023, первый заход поиска)
- Формат: онлайн, «ты компилятор» на Go + SQL с разбором EXPLAIN; заход начался с архитектурной задачи (см. `../system-design/`)

## 1. Указатель на строку

Что выведет код и почему?

```go
func setLinkHome(link *string) {
	*link = "http://home"
}

func main() {
	link := "http://other"

	setLinkHome(&link)

	fmt.Println(link)
}
```

## 2. panic и recover

Будет ли напечатан «ok»?

```go
func main() {
	defer func() {
		recover()
	}()

	panic("test panic")

	fmt.Println("ok")
}
```

## 3. Порядок вывода и done!

Функция должна напечатать:

```text
one
two
three
```

(в любом порядке), и в конце обязательно:

```text
done!
```

Но это не так. Исправьте код.

```go
func printText3(data []string) {
	for _, v := range data {
		go func() {
			fmt.Println(v)
		}()
	}

	fmt.Println("done!")
}

func main() {
	data := []string{"one", "two", "three"}

	printText3(data)
}
```

## 4. Параллельный счётчик

Мы пытаемся подсчитать количество выполненных параллельно операций. Что может пойти не так?

```go
var callCounter uint

func main() {
	wg := sync.WaitGroup{}

	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func() {
			// Ходим в базу, делаем долгую работу
			time.Sleep(time.Second)
			// Увеличиваем счетчик
			callCounter++
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("Call counter value = ", callCounter)
}
```

## 5. Таймаут обработки данных

Есть функция `processDataInternal`, которая может выполняться неопределённо долго. Чтобы контролировать процесс, мы добавили таймаут выполнения функции через context. Какие недостатки у кода ниже?

```go
type Service struct {
}

func (s *Service) processDataInternal(r io.Reader) error {
	return nil
}

func (s *Service) ProcessData(timeoutCtx context.Context, r io.Reader) error {
	errCh := make(chan error)

	go func() {
		errCh <- s.processDataInternal(r)
	}()

	select {
	case err := <-errCh:
		return err
	case <-timeoutCtx.Done():
		return timeoutCtx.Err()
	}
}
```

## 6. Что делает isCallAllowed

Опишите, что делает функция `isCallAllowed`:

```go
var callCount = make(map[uint]uint)

var locker = &sync.Mutex{}

func isCallAllowed(allowedCount uint) bool {
	if allowedCount == 0 {
		return true
	}

	locker.Lock()
	defer locker.Unlock()

	curTimeIndex := uint(time.Now().Unix() / 30)

	prevIndexVal, _ := callCount[curTimeIndex-1]

	if prevIndexVal >= allowedCount {
		return false
	}

	curIndexVal, ok := uint(0), false
	if curIndexVal, ok = callCount[curTimeIndex]; !ok {
		callCount[curTimeIndex] = 1

		return true
	}

	if (curIndexVal + prevIndexVal) >= allowedCount {
		return false
	}

	callCount[curTimeIndex]++

	return true
}
```

## 7. MySQL: пользователи без машины

У одного пользователя может быть неограниченное количество машин. Напишите запрос, который вернёт 10 пользователей, у которых нет авто с `car_id = 1`.

```text
+------------------+
|users             |
+------------------+
|id uint           | - PK
|name string       |
+------------------+

+------------------+
|user_cars         | - uniq index (user_id, car_id)
+------------------+
|user_id uint      |
|car_id uint       |
+------------------+

SELECT ...
```

## 8. MySQL: разбор EXPLAIN и оптимизация запроса

DevOps говорит, что в slowlog есть запрос, который выполняется дольше 10 секунд. Он отдал вам запрос, и вы вызвали EXPLAIN.

### 8.1. О чём говорит вывод EXPLAIN?

```text
+--+------------------+-----+----------+------+-------------------------------------------+-------------------------------------------+-------+-----------------+------+--------+------------------------+
|id|select_type       |table|partitions|type  |possible_keys                              |key                                        |key_len|ref               |rows  |filtered|Extra                   |
+--+------------------+-----+----------+------+-------------------------------------------+-------------------------------------------+-------+-----------------+------+--------+------------------------+
|1 |PRIMARY           |mc   |NULL      |ref   |idx_manager_id_client_id_uindex            |idx_manager_id_client_id_uindex            |1023   |const            |1     |100     |Using where; Using index|
|1 |PRIMARY           |m    |NULL      |eq_ref|idx_user_id                                |idx_user_id                                |1022   |bind.mc.client_id|1     |100     |Using where             |
|2 |DEPENDENT SUBQUERY|cdp  |NULL      |index |idx_client_id                              |idx_client_id                              |1022   |NULL             |189480|20.61   |Using where             |
+--+------------------+-----+----------+------+-------------------------------------------+-------------------------------------------+-------+-----------------+------+--------+------------------------+
```

### 8.2. Это сам запрос. Что можно сделать, чтобы он работал максимально быстро?

Задача запроса — выбрать клиентов определённого менеджера, у которых указано два этапа сделки:

```sql
SELECT m.*
FROM members m
         LEFT JOIN manager_clients mc on m.user_id = mc.client_id
WHERE mc.manager_id = '152734'
  AND m.user_id IN (SELECT client_id
                    FROM client_deal_phases cdp
                    WHERE cdp.phase_id IN (45, 47)
                    GROUP BY client_id
                    HAVING count(client_id) = 2);
```
