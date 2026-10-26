package main

import (
	"fmt"
)

func evaluateExpression(s string) int {
	if len(s) == 0 {
		return 0
	}

	sum := 0         // Итоговая сумма
	currProduct := 0 // Текущий результат умножения
	currNumber := 0  // Текущее число
	lastOp := '+'    // Последняя операция (по умолчанию '+')

	for i := 0; i < len(s); i++ {
		c := s[i]

		// Если символ - цифра, собираем число
		if c >= '0' && c <= '9' {
			currNumber = currNumber*10 + int(c-'0')
			continue
		}

		// Если символ - операция, обрабатываем предыдущее число
		if lastOp == '+' {
			sum += currProduct
			currProduct = currNumber
		} else if lastOp == '*' {
			currProduct *= currNumber
		}

		currNumber = 0
		lastOp = int32(c)
	}

	if lastOp == '+' {
		sum += currProduct + currNumber
	} else if lastOp == '*' {
		sum += currProduct * currNumber
	}

	return sum
}

func main() {
	tests := []string{
		"1",
		"1+2",
		"2*3",
		"1+2*3+5",
		"1*2+3*5*4",
		"1+2+3+4+5",
		"111+4+3*38+379*43",
	}

	for _, test := range tests {
		result := evaluateExpression(test)
		fmt.Printf("Expression: %s = %d\n", test, result)
	}
}
