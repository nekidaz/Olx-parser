package main

import (
	"Go-parser/parser"
	"fmt"
)

func main() {
	url := "https://www.olx.kz/list/q-хуй/"
	//url1 := "https://www.olx.kz/list/q-машина/"

	ads := parser.ParseAd(url)

	for _, ad := range ads {
		fmt.Println(ad.Link)

	}

}
