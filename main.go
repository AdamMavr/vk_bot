package main

import (
	"context"
	"database/sql"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"log"
	"lrn/Go/vk_bot/parser"
	"strings"
)

// TOKEN токен паблика
const TOKEN = "vk1.a.LFbHvvXgaj7mtfilwL7-1rt2_hrzH7JkvBAVBrEP4yDY_d1euPPLF0snfetL3gpOj7Cr-f0m6-sJvkyVihE7XlK8pzaketUW5EwzqyD9p9XqEar76W3TkavFfZx6B7F013SCTTp0_JLbWwcxnmuWyQt2_Mrgoxulgn6-oYnjnVUPRBmJZCSyRR0s8Nh1wjOxyiAw0o5avDee2KjdnNd31Q"

var (
	db         *sql.DB
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "password"
	dbName     = "chats"
)

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
	allHoroscopes, err := parser.GetTexts("https://t.me/s/neural_horo", ".tgme_widget_message_text")
	if err != nil {
		log.Println(err)
	}

	// New message event
	lp.MessageNew(func(_ context.Context, obj events.MessageNewObject) {
		log.Printf("%d: %s", obj.Message.ID, obj.Message.Text)

		if strings.HasPrefix(strings.ToLower(obj.Message.Text), "аска") {
			// Обрабатываем команду "гроб гороскоп знак_зодиака"
			if strings.Contains(obj.Message.Text, "гороскоп") {
				// Получаем знак зодиака из сообщения
				sign := strings.TrimSpace(strings.Split(obj.Message.Text, " ")[2])

				// Получаем гороскоп
				horoscope := parser.GetHoroscope(allHoroscopes, sign)

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
