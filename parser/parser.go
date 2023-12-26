package parser

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

// GetTexts возвращает текстовое представление элементов по selector на странице c переданным url
// Функция которая парсит все гороскопы в одну большую строку
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

	return strings.ToLower(strings.ReplaceAll(stringHoroscopes.String(), "\n", "")), nil
}

// GetHoroscope функция для получения конкретного гороскопа для знака зодиака
func GetHoroscope(allHoroscopes, sign string) string {
	strings.ToLower(sign)
	today := time.Now().Format("02.01.2006")
	start := strings.ToLower(sign)
	end := []string{"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"}

	startIndex := strings.Index(allHoroscopes, start)
	endIndex := strings.IndexAny(allHoroscopes[startIndex:], strings.Join(end, ""))

	if startIndex == -1 || endIndex == -1 {
		return ""
	}

	horoscope := allHoroscopes[startIndex+len(start) : startIndex+endIndex]

	// Возвращаем отформатированную строку
	return sign + " гороскоп на сегодня(" + today + ")" + horoscope
}
