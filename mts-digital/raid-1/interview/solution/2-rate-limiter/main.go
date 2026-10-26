package main

import (
	"fmt"
	"math"
	"time"
)

type Limiter interface {
	Allow(key string) bool
}

type tokenParams struct { // Про многопоточку: для конкурентного доступа нужен мьютекс
	max        int32
	speed      int32
	count      int32
	accessTime time.Time
}

type Baskets struct {
	basket map[string]tokenParams
}

func NewBaskets() Baskets {
	baskets := Baskets{
		basket: make(map[string]tokenParams),
	}

	return baskets
}

func (b Baskets) Allow(key string) bool {
	if _, exists := b.basket[key]; !exists {
		fmt.Printf("key %s not found", key)
		return false
	}

	creditedTokens := int32(0)
	now := time.Now()

	if !b.basket[key].accessTime.IsZero() {
		delta := now.Sub(b.basket[key].accessTime).Seconds()
		creditedTokens = int32(math.Floor(delta)) * b.basket[key].speed // за < 1 секунды начисляется 0

		if creditedTokens > b.basket[key].max {
			creditedTokens = b.basket[key].max
		}
	} else {
		creditedTokens = b.basket[key].max // первый раз начисляем максимум
	}

	newCount := b.basket[key].count - 1 + creditedTokens
	if newCount < 0 { // не снижаем токены, если токены потрачены
		return false
	}

	b.basket[key] = tokenParams{
		max:        b.basket[key].max,
		speed:      b.basket[key].speed,
		count:      newCount,
		accessTime: now,
	}

	return true
}

func main() {
	basket := NewBaskets()
	basket.basket["key"] = tokenParams{
		max:        5,
		speed:      1,
		count:      5,
		accessTime: time.Time{},
	}
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
	fmt.Println(basket.Allow("key"))
}
