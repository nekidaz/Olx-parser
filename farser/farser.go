package farser

import (
	"Go-parser/models"
	"Go-parser/utils"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
)

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

func ParseAd(url string) []models.AdModel {
	var ads []models.AdModel

	adCh := make(chan models.AdModel)
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

			parsedAd := models.AdModel{
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

	err := utils.ConvertToJson(ads)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON и записи в файл:", err)
	}
	return ads
}
