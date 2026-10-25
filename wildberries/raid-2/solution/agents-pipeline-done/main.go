package main

import (
	"fmt"
	"time"
)

type Agent struct {
	ID      int
	Enabled bool
}

func (a *Agent) Enable() { // Проблема 1: не применяются изменения к агенту
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
	done := make(chan struct{})
	go pipeEnableAndSend(pipe, agents)
	go pipeProcess(pipe, done)

	<-done // Проблема 2: нужно ожидать обработки агентов
	fmt.Println("final")
}

func addThirdPartyAgents(agents *[]Agent) {
	thirdParty := []Agent{
		{ID: 4},
		{ID: 5},
	}
	*agents = append(*agents, thirdParty...)
	// Проблема 3: надо передавать по указателю, иначе элементы не запишутся
}

func pipeEnableAndSend(pipe chan Enabler, agents []Agent) {
	for _, a := range agents {
		a.Enable()
		pipe <- &a // Проблема 4: теперь нужно передавать указатель потому что интерфейс a реализуется только с указателем
	}
	close(pipe) // Проблема 5: нужно закрыть канал, чтобы обработка прекратилась
}

func pipeProcess(pipe chan Enabler, done chan struct{}) {
	for a := range pipe { // Проблема 6: for range лучше чем for select и он работает пока канал не закрыт
		dbWrite(a)
	}
	close(done)
}

var dbWrite = func(a any) {
	fmt.Println(a)
	time.Sleep(time.Second * 1)
}
