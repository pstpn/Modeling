package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	defaultK   = 69069
	defaultC   = 7
	defaultMod = 4294967296
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Генератор случайных чисел")
	myWindow.Resize(fyne.Size{
		Width:  1080,
		Height: 500,
	})

	mainMenu := container.NewVBox()
	resultContainer := container.NewVBox()
	metricsLabel := widget.NewLabel("")

	itemsPerPage := 100

	methodLabel := widget.NewLabel("Выберите метод генерации:")
	methodOptions := widget.NewSelect([]string{"Табличный", "Алгоритмический"}, func(value string) {})
	methodOptions.Selected = "Табличный"

	countLabel := widget.NewLabel("Введите количество элементов:")
	countEntry := widget.NewEntry()
	countEntry.SetPlaceHolder("Введите число")

	generateButton := widget.NewButton("Сгенерировать", func() {
		currentPage := 0

		count := 0
		fmt.Sscanf(countEntry.Text, "%d", &count)
		if count < 1 {
			fyne.LogError("некорректное количество", fmt.Errorf("некорректное кол-во"))
			return
		}

		var sequence []float32
		if methodOptions.Selected == "Табличный" {
			sequence = generateTableSequence(count)
		} else {
			sequence = generateAlgorithmicSequence(count)
		}

		totalPages := (len(sequence) + itemsPerPage - 1) / itemsPerPage
		updateTable := func() {
			start := currentPage * itemsPerPage
			end := start + itemsPerPage
			if end > len(sequence) {
				end = len(sequence)
			}
			pageData := sequence[start:end]
			resultContainer.Objects = formatSequenceAsTable(pageData, 10)
			resultContainer.Refresh()
		}

		backButton := widget.NewButton("Вернуться в меню", func() {
			metricsLabel.Hide()
			myWindow.SetContent(mainMenu)
		})

		updateTable()
		metricsLabel.Show()
		metricsLabel.SetText(fmt.Sprintf("Метрики:\n\nСреднее: %.4f\nДисперсия: %.4f\nСтандартное отклонение: %.4f\nМин: %.4f\nМакс: %.4f\nРазмах: %.4f\n\n\n\nСтраница: %d из %d",
			mean(sequence), variance(sequence), stdDev(sequence), min(sequence), max(sequence), max(sequence)-min(sequence), currentPage+1, totalPages))

		myWindow.SetContent(container.NewBorder(
			nil, backButton, nil, metricsLabel, resultContainer,
		))

		myWindow.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
			if key.Name == fyne.KeyRight && currentPage < totalPages-1 {
				currentPage++
				updateTable()
			} else if key.Name == fyne.KeyLeft && currentPage > 0 {
				currentPage--
				updateTable()
			}

			metricsLabel.SetText(fmt.Sprintf("Метрики:\n\nСреднее: %.4f\nДисперсия: %.4f\nСтандартное отклонение: %.4f\nМин: %.4f\nМакс: %.4f\nРазмах: %.4f\n\n\n\nСтраница: %d из %d",
				mean(sequence), variance(sequence), stdDev(sequence), min(sequence), max(sequence), max(sequence)-min(sequence), currentPage+1, totalPages))
		})
	})

	mainMenu.Add(methodLabel)
	mainMenu.Add(methodOptions)
	mainMenu.Add(countLabel)
	mainMenu.Add(countEntry)
	mainMenu.Add(generateButton)
	mainMenu.Add(metricsLabel)

	myWindow.SetContent(mainMenu)

	myWindow.ShowAndRun()
}

func formatSequenceAsTable(sequence []float32, columns int) []fyne.CanvasObject {
	var labels []fyne.CanvasObject
	for _, num := range sequence {
		label := widget.NewLabel(fmt.Sprintf("%.4f", num))
		labels = append(labels, label)
	}
	return []fyne.CanvasObject{container.NewGridWithColumns(columns, labels...)}
}

func generateTableSequence(count int) []float32 {
	file, err := os.Open("/Users/stepa/Study/Modeling/lab_03/data/table_data.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil
	}
	defer file.Close()

	var sequence []float32
	scanner := bufio.NewScanner(file)
	for scanner.Scan() && len(sequence) < count {
		num, err := strconv.ParseFloat(scanner.Text(), 32)
		if err == nil {
			sequence = append(sequence, float32(num)) // Масштабируем в диапазон [0, 1]
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}

	return sequence
}

func generateAlgorithmicSequence(count int) []float32 {
	sequence := make([]float32, count)

	seed := uint32(time.Now().UnixNano() % int64(defaultMod))

	for i := 0; i < count; i++ {
		seed = uint32(int64(defaultK*seed+defaultC) % defaultMod)
		sequence[i] = float32(seed) / float32(defaultMod)
	}

	return sequence
}

func mean(sequence []float32) float32 {
	sum := float32(0)
	for _, num := range sequence {
		sum += num
	}
	return sum / float32(len(sequence))
}

func variance(sequence []float32) float32 {
	m := mean(sequence)
	var sum float32
	for _, num := range sequence {
		sum += (num - m) * (num - m)
	}
	return sum / float32(len(sequence))
}

func stdDev(sequence []float32) float32 {
	return sqrt(variance(sequence))
}

func sqrt(value float32) float32 {
	return float32(math.Sqrt(float64(value)))
}

func min(sequence []float32) float32 {
	m := sequence[0]
	for _, num := range sequence {
		if num < m {
			m = num
		}
	}
	return m
}

func max(sequence []float32) float32 {
	m := sequence[0]
	for _, num := range sequence {
		if num > m {
			m = num
		}
	}
	return m
}
