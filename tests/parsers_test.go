package tests

import (
	"Go-parser/farser"
	"Go-parser/parser"
	"testing"
)

func TestParserSpeed(t *testing.T) {
	url := "https://www.olx.kz/list/q-машина/"
	for i := 0; i < 5; i++ {
		parser.ParseAd(url)
	}

}

func TestFarserSpeed(t *testing.T) {
	url := "https://www.olx.kz/list/q-машина/"
	for i := 0; i < 5; i++ {
		farser.ParseAd(url)
	}
}
