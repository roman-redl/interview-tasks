package main

import (
	"fmt"
	"strconv"
)

func Eval(input string) int {
	if len(input) == 0 {
		return 0
	}

	res := 0    // Итоговая сумма
	num1 := 0   // Первое число или текущий результат умножения
	num2 := 0   // Второе число для умножения
	s := '+'    // Предыдущая операция
	numInd := 0 // Индекс начала текущего числа

	for i := 0; i < len(input); i++ {
		c := input[i]

		if c == '+' || c == '*' {
			num, _ := strconv.Atoi(input[numInd:i])
			if num1 == 0 {
				num1 = num
			} else {
				num2 = num
				if s == '*' {
					num1 *= num2
					num2 = 0
				} else if s == '+' {
					res += num1
					num1 = num2
					num2 = 0
				}
			}
			s = int32(c) // Запоминаем операцию
			numInd = i + 1
		}
	}

	// Обрабатываем последнее число
	num, _ := strconv.Atoi(input[numInd:])
	if s == '+' {
		if num1 != 0 {
			res += num1
			num1 = num
		} else {
			num1 = num
		}
	} else if s == '*' {
		if num1 == 0 {
			num1 = num
		} else {
			num1 *= num
		}
	}

	res += num1

	return res
}

func main() {
	tests := []string{
		"1",
		"1+2",
		"2*3",
		"1+2*30+5",
		"1*2+3*5*4",
		"1+2+3+4+5",
		"111+4+3*38+379*43",
	}

	for _, test := range tests {
		result := Eval(test)
		fmt.Printf("Expression: %s = %d\n", test, result)
	}
}
