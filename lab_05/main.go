package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Operator struct {
	ID      int
	MinTime int
	MaxTime int
	IsBusy  bool
	Mutex   *sync.Mutex
}

type Computer struct {
	ID          int
	ProcessTime int
	Mutex       *sync.Mutex
}

func simulateRequests(totalClients int, updateStatus func(string)) float64 {
	wg := &sync.WaitGroup{}
	opMutex := &sync.Mutex{}
	compMutex := &sync.Mutex{}
	mutex := &sync.Mutex{}

	operators := []Operator{
		{ID: 1, MinTime: 150, MaxTime: 250, Mutex: opMutex},
		{ID: 2, MinTime: 300, MaxTime: 500, Mutex: opMutex},
		{ID: 3, MinTime: 200, MaxTime: 600, Mutex: opMutex},
	}

	computers := []Computer{
		{ID: 1, ProcessTime: 150, Mutex: compMutex},
		{ID: 2, ProcessTime: 300, Mutex: compMutex},
	}

	var servedClients int
	var rejectedClients int
	allTime := 0 * time.Nanosecond

	for i := 0; i < totalClients; i++ {
		i := i + 1
		t := allTime
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			selectedOperator := -1
			for j, op := range operators {
				op.Mutex.Lock()
				if !op.IsBusy {
					selectedOperator = j
					operators[j].IsBusy = true
					op.Mutex.Unlock()
					break
				}
				op.Mutex.Unlock()
			}

			if selectedOperator == -1 {
				mutex.Lock()
				rejectedClients++
				mutex.Unlock()
				updateStatus(fmt.Sprintf("Клиент %d получил отказ в %d минут(ы)", clientID, t.Milliseconds()/10))
				return
			}

			op := &operators[selectedOperator]
			processTime := time.Duration(rand.Intn(op.MaxTime-op.MinTime+1)+op.MinTime) * time.Millisecond
			t += processTime
			time.Sleep(processTime)
			updateStatus(fmt.Sprintf("Клиент %d был обработан оператором в %d минут(ы)", clientID, t.Milliseconds()/10))

			op.Mutex.Lock()
			op.IsBusy = false
			op.Mutex.Unlock()

			computerID := 0
			if op.ID == 3 {
				computerID = 1
			}

			computer := &computers[computerID]
			computer.Mutex.Lock()
			compProcessTime := time.Duration(computer.ProcessTime) * time.Millisecond
			t += compProcessTime
			time.Sleep(compProcessTime)
			computer.Mutex.Unlock()
			updateStatus(fmt.Sprintf("Клиент %d был обработан компьютером в %d минут(ы)", clientID, t.Milliseconds()/10))

			mutex.Lock()
			servedClients++
			mutex.Unlock()
			updateStatus(fmt.Sprintf("Клиент %d обслужен в %d минут(ы)", clientID, t.Milliseconds()/10))
		}(i)

		sleepTime := time.Duration(rand.Intn(40)+80) * time.Millisecond
		allTime += sleepTime
		time.Sleep(sleepTime)
	}

	wg.Wait()

	totalClientsProcessed := servedClients + rejectedClients
	if totalClientsProcessed == 0 {
		return 0
	}
	return float64(rejectedClients) / float64(totalClientsProcessed)
}

func main() {
	a := app.New()
	w := a.NewWindow("Моделирование информационного центра")

	messageList := container.NewVBox()
	scroll := container.NewVScroll(messageList)
	scroll.SetMinSize(fyne.NewSize(600, 400))

	updateStatus := func(msg string) {
		messageList.Add(widget.NewLabel(msg))
		scroll.ScrollToBottom()
	}

	startButton := widget.NewButton("Запустить моделирование", func() {
		messageList.RemoveAll()
		messageList.Refresh()
		go func() {
			probability := simulateRequests(300, updateStatus)
			updateStatus(fmt.Sprintf("\nВероятность отказа: %.2f%%", probability*100))
		}()
	})

	w.SetContent(container.NewVBox(
		widget.NewLabel("Моделирование информационного центра"),
		scroll,
		layout.NewSpacer(),
		startButton,
	))

	w.Resize(fyne.NewSize(800, 500))
	w.ShowAndRun()
}
