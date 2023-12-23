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

	fmt.Println(allHoroscopes)
	//fmt.Println(GetHoroscope(AllHoroscopes, "Телец"))

}

// GetTexts возвращает текстовое представление элементов по selector на странице c переданным url
func GetTexts(url, selector string) (string, error) {

	// Скачиваем html-страницу
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	// Считываем страницу в goquery-документ
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	// Находим все элементы с переданным селектором и сохраняем их содержимое в строку
	var stringHoroscopes strings.Builder
	doc.Find(selector).Last().Each(func(i int, s *goquery.Selection) {
		stringHoroscopes.WriteString(s.Text())
	})

	return strings.ReplaceAll(stringHoroscopes.String(), "\n", ""), nil
}

func GetHoroscope(allHoroscopes, sign string) string {
	start := sign
	end := []string{"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"}

	startIndex := strings.Index(allHoroscopes, start)
	endIndex := strings.IndexAny(allHoroscopes[startIndex:], strings.Join(end, ""))

	if startIndex == -1 || endIndex == -1 {
		return ""
	}

	return sign + allHoroscopes[startIndex+len(start):startIndex+endIndex]
}
