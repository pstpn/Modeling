package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type UniformDistribution struct {
	a, b float64
}

func (ud *UniformDistribution) Generate() float64 {
	return ud.a + rand.Float64()*(ud.b-ud.a)
}

type NormalDistribution struct {
	mean, stddev float64
}

func (nd *NormalDistribution) Generate() float64 {
	return rand.NormFloat64()*nd.stddev + nd.mean
}

func stepModel(generator *UniformDistribution, processor *NormalDistribution, count int, repeatProb float64, step float64) int {
	tasksDone := 0
	curQueueLen := 0
	maxQueueLen := 0
	countGen := 0
	timeCurrent := step
	timeGenerated := generator.Generate()
	timeProcessed := 0.0

	for countGen < count {
		if timeCurrent > timeGenerated {
			curQueueLen++
			if curQueueLen > maxQueueLen {
				maxQueueLen = curQueueLen
			}
			timeGenerated += generator.Generate()
			countGen++
		}

		if timeCurrent > timeProcessed {
			if curQueueLen > 0 {
				curQueueLen--
				tasksDone++
				if rand.Float64()*100 <= repeatProb {
					curQueueLen++
				}
				timeProcessed += processor.Generate()
			}
		}

		timeCurrent += step
	}
	return maxQueueLen
}

func eventModel(generator *UniformDistribution, processor *NormalDistribution, count int, repeatProb float64) int {
	tasksDone := 0
	curQueueLen := 0
	maxQueueLen := 0
	countGen := 1
	free := true
	events := []struct {
		time  float64
		event string
	}{{time: generator.Generate(), event: "generate"}}

	for countGen < count {
		event := events[0]
		events = events[1:]

		switch event.event {
		case "generate":
			curQueueLen++
			if curQueueLen > maxQueueLen {
				maxQueueLen = curQueueLen
			}
			events = append(events, struct {
				time  float64
				event string
			}{time: event.time + generator.Generate(), event: "generate"})
			countGen++

			if free {
				events = append(events, struct {
					time  float64
					event string
				}{time: event.time + processor.Generate(), event: "process"})
				free = false
			}

		case "process":
			tasksDone++
			if rand.Float64()*100 <= repeatProb {
				curQueueLen++
			}
			if curQueueLen > 0 {
				curQueueLen--
				events = append(events, struct {
					time  float64
					event string
				}{time: event.time + processor.Generate(), event: "process"})
			} else {
				free = true
			}
		}
	}
	return maxQueueLen
}

func main() {
	a := app.New()
	w := a.NewWindow("Система массового обслуживания")

	w.Resize(fyne.NewSize(500, 500))

	aEntry := widget.NewEntry()
	bEntry := widget.NewEntry()
	meanEntry := widget.NewEntry()
	stddevEntry := widget.NewEntry()
	countEntry := widget.NewEntry()
	repeatEntry := widget.NewEntry()
	stepEntry := widget.NewEntry()

	resultStep := widget.NewLabel("Размер очереди (Пошаговый подход):")
	resultEvent := widget.NewLabel("Размер очереди (Событийный подход):")

	runButton := widget.NewButton("Запустить симуляцию", func() {
		a, _ := strconv.ParseFloat(aEntry.Text, 64)
		b, _ := strconv.ParseFloat(bEntry.Text, 64)
		mean, _ := strconv.ParseFloat(meanEntry.Text, 64)
		stddev, _ := strconv.ParseFloat(stddevEntry.Text, 64)
		count, _ := strconv.Atoi(countEntry.Text)
		repeat, _ := strconv.ParseFloat(repeatEntry.Text, 64)
		step, _ := strconv.ParseFloat(stepEntry.Text, 64)

		generator := &UniformDistribution{a: a, b: b}
		processor := &NormalDistribution{mean: mean, stddev: stddev}

		stepResult := stepModel(generator, processor, count, repeat, step)
		eventResult := eventModel(generator, processor, count, repeat)

		resultStep.SetText(fmt.Sprintf("Размер очереди (Пошаговый подход): %d", stepResult))
		resultEvent.SetText(fmt.Sprintf("Размер очереди (Событийный подход): %d", eventResult))
	})

	form := container.NewVBox(
		widget.NewLabel("Генератор (Равномерное распределение)"),
		widget.NewForm(widget.NewFormItem("a", aEntry), widget.NewFormItem("b", bEntry)),
		widget.NewLabel("Обработчик (Нормальное распределение)"),
		widget.NewForm(widget.NewFormItem("Мат ожидание", meanEntry), widget.NewFormItem("Среднекв отклонение", stddevEntry)),
		widget.NewForm(widget.NewFormItem("Кол-во заявок", countEntry), widget.NewFormItem("Вероятность возврата (%)", repeatEntry), widget.NewFormItem("Шаг (с)", stepEntry)),
		runButton,
		resultStep,
		resultEvent,
	)

	w.SetContent(form)
	w.ShowAndRun()
}
