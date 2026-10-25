package main

import (
	"fmt"
	"sync"
	"time"
)

type Agent struct {
	ID      int
	Enabled bool
}

func (a *Agent) Enable() {
	a.Enabled = true
}

type Enabler interface {
	Enable()
}

// 1. Инициализация слайса агентов
// 2. Дополнение слайса сторонними агентами
// 3. Потоковая обработка объектов - активация и отправка на выполнение работы
// 4. Потоковая обработка объектов - сохранение в БД и распечатка результатов
func main() {
	agents := make([]Agent, 0, 5)
	for i := 0; i < 2; i++ {
		agents = append(agents, Agent{ID: i})
	}

	addThirdPartyAgents(&agents)

	pipe := make(chan Enabler)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		pipeEnableAndSend(pipe, agents)
		close(pipe)
	}()

	go func() {
		defer wg.Done()
		pipeProcess(pipe)
	}()

	wg.Wait()
}

func addThirdPartyAgents(agents *[]Agent) {
	thirdParty := []Agent{
		{ID: 4},
		{ID: 5},
	}
	*agents = append(*agents, thirdParty...)
}

func pipeEnableAndSend(pipe chan Enabler, agents []Agent) {
	for _, a := range agents {
		a.Enable()
		pipe <- &a
	}
}

func pipeProcess(pipe chan Enabler) {
	for a := range pipe {
		dbWrite(a)
	}
}

var dbWrite = func(a any) {
	fmt.Println(a)
	time.Sleep(time.Second * 1)
}
