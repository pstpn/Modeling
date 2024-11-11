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
		Width:  800,
		Height: 600,
	})

	resultContainer := container.NewVBox()
	metricsLabel := widget.NewLabel("")

	userInputEntries := createUserInputTable()

	generateButton := widget.NewButton("Сгенерировать", func() {
		sequenceTable1 := generateTableSequence(10, 1)
		sequenceTable2 := generateTableSequence(10, 2)
		sequenceTable3 := generateTableSequence(10, 3)

		sequenceAlg1 := generateAlgorithmicSequence(10, 1)
		sequenceAlg2 := generateAlgorithmicSequence(10, 2)
		sequenceAlg3 := generateAlgorithmicSequence(10, 3)

		userSequence := getUserInput(userInputEntries)

		tableResults := formatSequenceAsTable(sequenceTable1, sequenceTable2, sequenceTable3, "Табличный метод")
		algResults := formatSequenceAsTable(sequenceAlg1, sequenceAlg2, sequenceAlg3, "Алгоритмический метод")

		resultContainer.Objects = []fyne.CanvasObject{tableResults, algResults}
		resultContainer.Refresh()

		metricsLabel.SetText(fmt.Sprintf("Метрика колебаний последовательности:\n\nТабличный метод:\n0-9:     %.4f\n0-99:   %.4f\n0-999: %.4f\n\nАлгоритмический метод:\n0-9:     %.4f\n0-99:   %.4f\n0-999: %.4f\n\nПользовательская таблица:\n\n %.4f",
			calculateFluctuationMetric(sequenceTable1), calculateFluctuationMetric(sequenceTable2), calculateFluctuationMetric(sequenceTable3),
			calculateFluctuationMetric(sequenceAlg1), calculateFluctuationMetric(sequenceAlg2), calculateFluctuationMetric(sequenceAlg3),
			calculateFluctuationMetric(userSequence)))
	})

	myWindow.SetContent(container.NewVBox(
		container.NewHBox(resultContainer, metricsLabel),
		container.NewVBox(widget.NewLabel("Пользовательская таблица:"), container.NewGridWithColumns(10, convertToCanvasObjects(userInputEntries)...)),
		generateButton,
	))

	myWindow.ShowAndRun()
}

func convertToCanvasObjects(entries []*widget.Entry) []fyne.CanvasObject {
	var objects []fyne.CanvasObject
	for _, entry := range entries {
		objects = append(objects, entry)
	}
	return objects
}

func createUserInputTable() []*widget.Entry {
	entries := make([]*widget.Entry, 10)
	for i := 0; i < 10; i++ {
		entry := widget.NewEntry()
		entry.SetText("0")
		entries[i] = entry
	}
	return entries
}

func getUserInput(entries []*widget.Entry) []int {
	var values []int
	for _, entry := range entries {
		val, err := strconv.Atoi(entry.Text)
		if err != nil {
			val = 0
		}
		values = append(values, val)
	}
	return values
}

func formatSequenceAsTable(seq1, seq2, seq3 []int, title string) *fyne.Container {
	label := widget.NewLabel(title)
	table1 := container.NewGridWithColumns(10, formatAsLabels(seq1)...)
	table2 := container.NewGridWithColumns(10, formatAsLabels(seq2)...)
	table3 := container.NewGridWithColumns(10, formatAsLabels(seq3)...)
	return container.NewVBox(label, table1, table2, table3)
}

func formatAsLabels(sequence []int) []fyne.CanvasObject {
	var labels []fyne.CanvasObject
	for _, num := range sequence {
		labels = append(labels, widget.NewLabel(fmt.Sprintf("%d", num)))
	}
	return labels
}

func generateTableSequence(count int, n int) []int {
	file, err := os.Open("/Users/stepa/Study/Modeling/lab_03/data/table_data.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil
	}
	defer file.Close()
	posFile, err := os.Open("/Users/stepa/Study/Modeling/lab_03/data/pos.txt")
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil
	}
	defer posFile.Close()

	var index int
	_, err = fmt.Fscanf(posFile, "%d", &index)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return nil
	}

	var sequence []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		num, err := strconv.ParseFloat(scanner.Text(), 32)
		if err == nil {
			sequence = append(sequence, int(num*math.Pow10(n)))
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
	}

	res := make([]int, count)
	for i := range count {
		if index+6 >= len(sequence) {
			index %= 100
		}
		index += 6

		res[i] = sequence[index]
	}

	err = os.WriteFile("/Users/stepa/Study/Modeling/lab_03/data/pos.txt", []byte(strconv.Itoa(index)), 0777)
	if err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}

	return res
}

func generateAlgorithmicSequence(count, n int) []int {
	sequence := make([]int, count)

	seed := uint32(time.Now().UnixNano() % int64(defaultMod))
	for i := 0; i < count; i++ {
		seed = uint32(int64(defaultK*seed+defaultC) % defaultMod)
		sequence[i] = int(float32(seed) / float32(defaultMod) * float32(math.Pow10(n)))
	}

	return sequence
}

func calculateFluctuationMetric(sequence []int) float64 {
	var differences []float64
	for i := 1; i < len(sequence); i++ {
		diff := math.Abs(float64(sequence[i] - sequence[i-1]))
		differences = append(differences, diff)
	}
	meanDiff := mean(differences)

	absDeviations := make([]float64, len(differences))
	for i, diff := range differences {
		absDeviations[i] = math.Abs(diff - meanDiff)
	}
	meanAbsDeviation := mean(absDeviations)

	if meanDiff == 0 {
		return 0
	}
	metric := meanAbsDeviation / meanDiff
	return metric
}

func mean(sequence []float64) float64 {
	sum := 0.0
	for _, num := range sequence {
		sum += num
	}
	return sum / float64(len(sequence))
}
