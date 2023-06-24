package parser

import (
	"Go-parser/models"
	"Go-parser/utils"
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
)

func parseTitle(adElement soup.Root) string {
	titleElement := adElement.Find("h6")
	if titleElement.Error != nil {
		return ""
	}
	return titleElement.Text()
}

func parseLocation(adElement soup.Root) string {
	locationElement := adElement.Find("p", "data-testid", "location-date")
	if locationElement.Error != nil {
		return ""
	}
	return locationElement.Text()
}

func parsePrice(adElement soup.Root) string {
	priceElement := adElement.Find("p", "data-testid", "ad-price")
	if priceElement.Error != nil {
		return ""
	}
	return priceElement.Text()
}

func parseCondition(adElement soup.Root) string {
	conditionElement := adElement.Find("span").Find("span")
	if conditionElement.Error != nil {
		return ""
	}
	return conditionElement.Attrs()["title"]
}

func parseLink(adElement soup.Root) string {
	linkElement := adElement.Find("a").Attrs()["href"]
	fullLink := fmt.Sprintf("https://www.olx.kz/%s", linkElement)
	return fullLink
}

func ParseAd(url string) {
	var ads []models.AdModel
	for {
		resp, err := soup.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		doc := soup.HTMLParse(resp)
		adElements := doc.FindAll("div", "data-cy", "l-card")

		for _, ad := range adElements {
			title := parseTitle(ad)
			price := parsePrice(ad)
			location := parseLocation(ad)
			condition := parseCondition(ad)
			link := parseLink(ad)

			parsedAd := models.AdModel{
				Title:     title,
				Price:     price,
				Location:  location,
				Condition: condition,
				Link:      link,
			}

			ads = append(ads, parsedAd)
		}

		nextPageLink := doc.Find("a", "data-testid", "pagination-forward")
		if nextPageLink.Error == nil {
			nextPage := fmt.Sprintf("https://www.olx.kz%s", nextPageLink.Attrs()["href"])
			url = nextPage
		} else {
			break
		}
	}

	err := utils.ConvertToJson(ads)
	if err != nil {
		fmt.Println("Ошибка при преобразовании в JSON и записи в файл:", err)
	}
}
