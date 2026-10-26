package main

import "fmt"

// isMonotonic: слайс монотонен, если он целиком не убывает или целиком
// не возрастает. Равные соседние элементы разрешены в обе стороны.
func isMonotonic(a []int) bool {
	inc, dec := true, true
	for i := 1; i < len(a); i++ {
		switch {
		case a[i] > a[i-1]:
			dec = false
		case a[i] < a[i-1]:
			inc = false
		}
		if !inc && !dec {
			return false
		}
	}
	return inc || dec
}

func main() {
	tests := []struct {
		in   []int
		want bool
	}{
		{[]int{1, 7}, true},
		{[]int{1, 1}, true},
		{[]int{3, 3, 1}, true},
		{[]int{9, 5, 1}, true},
		{[]int{23, 5, 23}, false},
		{[]int{5}, true},
		{[]int{1, 2, 3, 2}, false},
	}
	for _, t := range tests {
		got := isMonotonic(t.in)
		status := "OK "
		if got != t.want {
			status = "FAIL"
		}
		fmt.Printf("%s isMonotonic(%v) = %v\n", status, t.in, got)
	}
}
