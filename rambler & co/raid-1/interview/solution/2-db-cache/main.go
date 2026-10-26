// Есть функция получения данных из "базы".
// Необходимо придумать вариант как нам правильно читать данные из бд
// Данные мы запрашиваем часто и много, тут для примера используем просто цикл, который пытается получить данные

package main

import (
	"fmt"
	"sync"
	"time"
)

type User struct {
	ID   int
	Name string
}

var (
	dbHits    int
	cacheHits int
)

var (
	users = make(map[int]User)
	mu    = sync.Mutex{}
)

func getUserInfoFromDB(id int) User {
	time.Sleep(100 * time.Millisecond)
	dbHits++
	return User{ID: id, Name: fmt.Sprintf("User-%d", id)}
}

func getUserInfoFromCache(id int) (User, bool) {
	mu.Lock()
	defer mu.Unlock()
	user, ok := users[id]
	if !ok {
		return user, ok
	}
	return user, ok
}

func setUserInfoToCache(user User) {
	mu.Lock()
	defer mu.Unlock()
	users[user.ID] = user
}

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		time.Sleep(200 * time.Millisecond) // без задержки будет 5 db hits
		go func(i int) {
			defer wg.Done()
			user, ok := getUserInfoFromCache(i)
			if !ok {
				user = getUserInfoFromDB(i)
				setUserInfoToCache(user)
			}
			fmt.Println(user)
		}(42)
	}
	wg.Wait()
	fmt.Println("DB hits:", dbHits)
}
