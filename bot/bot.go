package main

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"net/http"
	"strings"
)

// TOKEN токен паблика
const TOKEN = "vk1.a.LFbHvvXgaj7mtfilwL7-1rt2_hrzH7JkvBAVBrEP4yDY_d1euPPLF0snfetL3gpOj7Cr-f0m6-sJvkyVihE7XlK8pzaketUW5EwzqyD9p9XqEar76W3TkavFfZx6B7F013SCTTp0_JLbWwcxnmuWyQt2_Mrgoxulgn6-oYnjnVUPRBmJZCSyRR0s8Nh1wjOxyiAw0o5avDee2KjdnNd31Q"

// стартовая точка программы
func main() {
	vk := api.NewVK(TOKEN)

	// получаем информацию о группе (id через токен)
	group, err := vk.GroupsGetByID(nil)
	if err != nil {
		log.Fatal(err)
	}

	// инициализация лонг пула
	lp, err := longpoll.NewLongPoll(vk, group[0].ID)
	if err != nil {
		log.Fatal(err)
	}

	// парсим гороскопы из телеграма по селектору css
	allHoroscopes, err := GetTexts("https://t.me/s/neural_horo", ".tgme_widget_message_text")
	if err != nil {
		log.Println(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.ID, obj.Message.Text)

		if strings.HasPrefix(obj.Message.Text, "гроб") {
			// Обрабатываем команду "гроб гороскоп знак_зодиака"
			if strings.Contains(obj.Message.Text, "гороскоп") {
				// Получаем знак зодиака из сообщения
				sign := strings.TrimSpace(strings.Split(obj.Message.Text, " ")[2])

				// Получаем гороскоп
				horoscope := GetHoroscope(allHoroscopes, sign)

				// Отправляем гороскоп в ответ
				b := params.NewMessagesSendBuilder()
				b.Message(horoscope)
				b.RandomID(0)
				b.PeerID(obj.Message.PeerID)

				_, err = vk.MessagesSend(b.Params)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	})

	// запускаем бота
	log.Println("Start Long Poll")
	if err := lp.Run(); err != nil {
		log.Fatal(err)
	}
}

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

	return strings.ReplaceAll(stringHoroscopes.String(), "\n", ""), nil
}

// GetHoroscope функция для получения конкретного гороскопа для знака зодиака
func GetHoroscope(allHoroscopes, sign string) string {
	start := sign
	end := []string{"♈", "♉", "♊", "♋", "♌", "♍", "♎", "♏", "♐", "♑", "♒", "♓"}

	startIndex := strings.Index(allHoroscopes, start)
	endIndex := strings.IndexAny(allHoroscopes[startIndex:], strings.Join(end, ""))

	if startIndex == -1 || endIndex == -1 {
		return ""
	}

	horoscope := allHoroscopes[startIndex+len(start) : startIndex+endIndex]

	// Возвращаем отформатированную строку
	return sign + horoscope
}
