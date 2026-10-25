# Wildberries — заход 7, техническое интервью

- Даты: 2025-11-24
- Формат: «что произойдёт / что выведется»

## 1. Изменение символа строки

Что произойдёт?

```go
func main() {
	s := "asdq"
	s[0] = 'r'
	println(s)
}
```

## 2. Слайсы и append

Что выведется?

```go
func slices() {
	a := []string{"a", "b", "c"}
	b := a[1:2]
	b[0] = "q"
	b = append(b, "d")
	b = append(b, "e")
	b[0] = "x"
	fmt.Println(a)
}
```

## 3. range по строке с кириллицей

Что выведется?

```go
func printString() {
	s := "Юadsf"

	for i := range s {
		fmt.Println(i)
	}
}
```
