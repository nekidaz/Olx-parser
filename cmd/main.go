package main

import (
	"Go-parser/farser"
)

func main() {
	url := "https://www.olx.kz/list/q-машина/"

	farser.ParseAd(url)

	// Вывод результатов
}
