package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func gen() {
	file, err := os.Create("data/table_data.txt")
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000000; i++ {
		_, err := file.WriteString(fmt.Sprintf("%f\n", rand.Float32()))
		if err != nil {
			fmt.Println("Ошибка при записи в файл:", err)
			return
		}
	}

	fmt.Println("Файл с некоррелированными числами успешно создан.")
}
