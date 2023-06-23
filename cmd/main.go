package main

import (
	"Go-parser/farser"
	"Go-parser/parser"
	"fmt"
	"time"
)

func main() {
	url := "https://www.olx.kz/list/q-какашки/"

	// Замер времени выполнения парсера "farser"
	startTime := time.Now()
	farser.ParseAd(url)
	endTime := time.Now()
	farserExecutionTime := endTime.Sub(startTime)

	// Замер времени выполнения парсера "parser"
	startTime = time.Now()
	parser.ParseAd(url)
	endTime = time.Now()
	parserExecutionTime := endTime.Sub(startTime)

	// Вывод результатов
	fmt.Println("Время выполнения парсера 'farser':", farserExecutionTime)
	fmt.Println("Время выполнения парсера 'parser':", parserExecutionTime)
}
