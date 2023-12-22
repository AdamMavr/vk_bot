package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	allHoroscopes, err := GetTexts("https://t.me/s/neural_horo", ".tgme_widget_message_text")
	if err != nil {
		log.Println(err)
	}

	stringHoroscopes := strings.Join(allHoroscopes, "\n")

	fmt.Println(stringHoroscopes)

	for _, t := range allHoroscopes {
		stringHoroscopes += t
	}
}

// GetTexts возвращает текстовое представление элементов
// по selector на странице c переданным url
func GetTexts(url, selector string) ([]string, error) {

	// Скачиваем html-страницу
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	// Считываем страницу в goquery-документ
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// Находим все элементы с переданным селектором
	// и сохраняем их содержимое в список
	var texts []string
	doc.Find(selector).First().Each(func(i int, s *goquery.Selection) {
		texts = append(texts, s.Text())
	})

	return texts, nil
}
