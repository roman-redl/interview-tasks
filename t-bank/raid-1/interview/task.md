# Т-Банк — заход 1, техническое интервью

- Даты: 2025-07-29
- Формат: онлайн, задачи «что выведется» и на написание функции + теория

Обвязка package/import в задачах опущена. Четвёртая, самая объёмная задача (код-ревью большого «джунского» фрагмента с последующим траблшутингом уже в продакшене) была описана интервьюером, но не разбиралась — к её началу закончилось время; взамен досконально проговорили теорию.

## 1. Append со слайсом указателей: что выведется

До запуска определить, что будет выведено в output. После ответа — дополнительные вопросы: из чего состоит слайс под капотом; почему длина не изменилась снаружи; как бы ты сделал, чтобы изменение возымело эффект; что изменится, если range идёт не по numbers, а по срезу numbers, и почему при выводе до пятого элемента всё ломается, а до четвёртого — нет.

```go
func appendLen(numbers []*int) {
	size := len(numbers)
	numbers = append(numbers, &size)
}

func main() {
	numbers := make([]*int, 0, 5)
	var number int
	for range 3 {
		number++
		numbers = append(numbers, &number)
	}
	appendLen(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number)
	}
}
```

## 2. Функция double

В тот же код после `appendLen` добавляется вызов `double(numbers)`. Реализовать функцию ровно так, как описано в комментарии над ней. После решения — дополнительный вопрос: у функции появился сайд-эффект (меняются значения по указателям, которые могут принадлежать внешней программе) — попробовать реализовать, не заменяя указатели в слайсе.

```go
	double(numbers)

	for _, number := range numbers {
		fmt.Printf("%d ", *number)
	}
}

// 3 3 3 -> 6 6 6
func double(numbers []*int) {

}
```

## 3. RunWriter: что происходит в программе

Сказать, как программа себя поведёт и что окажется на экране. Задача состоит из двух частей; вторая (закомментированный вариант main с `RunProcessor` и `sync.WaitGroup`) не разбиралась — закончилось время.

```go
func RunWriter() <-chan map[string]int {
	var prices = make(chan map[string]int)
	go func() {
		var currentPrice = map[string]int{
			"AAPL": 163,
			"USD":  117,
			"EUR":  124,
			"NVDA": 234,
		}
		for i := 1; i < 5; i++ {
			for key, value := range currentPrice {
				currentPrice[key] = int(float64(value) * 1.3)
			}
			prices <- currentPrice
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return prices
}

func main() {
	p := RunWriter()
	var prices []map[string]int

	for {
		prices = append(prices, <-p)
	}

	for _, price := range prices {
		fmt.Println(price)
	}
}
```
