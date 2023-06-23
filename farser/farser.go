package farser

import (
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"log"
	"os"
)

type AdParse struct {
	Title     string `json:"title"`
	Price     string `json:"price"`
	Location  string `json:"location"`
	Condition string `json:"condition"`
	Link      string `json:"link"`
}

func parseTitle(adElement soup.Root, ch chan string) {
	titleElement := adElement.Find("h6").Text()
	ch <- titleElement
}

func parseLocation(adElement soup.Root, ch chan string) {
	locationElement := adElement.Find("p", "data-testid", "location-date")
	ch <- locationElement.Text()
}

func parsePrice(adElement soup.Root, ch chan string) {
	priceElement := adElement.Find("p", "data-testid", "ad-price")
	if priceElement.Error != nil {
		ch <- ""
	} else {
		ch <- priceElement.Text()
	}
}

func parseCondition(adElement soup.Root, ch chan string) {
	conditionElement := adElement.Find("span").Find("span")
	if conditionElement.Error != nil {
		ch <- ""
	} else {
		ch <- conditionElement.Attrs()["title"]
	}
}

func parseLink(adElement soup.Root, ch chan string) {
	linkElement := adElement.Find("a").Attrs()["href"]
	fullLink := fmt.Sprintf("https://www.olx.kz/%s", linkElement)
	ch <- fullLink
}

func ParseAd(url string) []AdParse {
	var ads []AdParse
	adCh := make(chan AdParse)
	doneCh := make(chan struct{})

	go func() {
		for ad := range adCh {
			ads = append(ads, ad)
		}
		doneCh <- struct{}{}
	}()

	for {
		resp, err := soup.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		doc := soup.HTMLParse(resp)
		adElements := doc.FindAll("div", "data-cy", "l-card")

		for _, ad := range adElements {
			titleCh := make(chan string)
			priceCh := make(chan string)
			locationCh := make(chan string)
			conditionCh := make(chan string)
			linkCh := make(chan string)

			go parseTitle(ad, titleCh)
			go parsePrice(ad, priceCh)
			go parseLocation(ad, locationCh)
			go parseCondition(ad, conditionCh)
			go parseLink(ad, linkCh)

			parsedAd := AdParse{
				Title:     <-titleCh,
				Price:     <-priceCh,
				Location:  <-locationCh,
				Condition: <-conditionCh,
				Link:      <-linkCh,
			}

			adCh <- parsedAd

			close(titleCh)
			close(priceCh)
			close(locationCh)
			close(conditionCh)
			close(linkCh)
		}

		nextPageLink := doc.Find("a", "data-testid", "pagination-forward")
		if nextPageLink.Error == nil {
			nextPage := fmt.Sprintf("https://www.olx.kz%s", nextPageLink.Attrs()["href"])
			url = nextPage
		} else {
			close(adCh)
			break
		}
	}

	<-doneCh

	err := convertToJson(ads)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON и записи в файл:", err)
	}
	return ads
}

func convertToJson(ads []AdParse) error {
	jsonData, err := json.MarshalIndent(ads, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll("data", 0755) // Создание директории, если она не существует
	if err != nil {
		return err
	}

	filePath := "data/data1.json"
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Данные успешно записаны в файл", filePath)
	return nil
}
