package utils

import (
	"Go-parser/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func ConvertToJson(ads []models.AdModel) error {
	jsonData, err := json.MarshalIndent(ads, "", "  ")
	if err != nil {
		return err
	}

	// Генерация уникального имени файла на основе текущей даты и времени
	currentTime := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("data/data_%s.json", currentTime)

	err = os.MkdirAll("data", 0755) // Создание директории, если она не существует
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("Данные успешно записаны в файл", fileName)
	return nil
}
