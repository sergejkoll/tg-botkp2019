//package main
//
//import (
//	"log"
//	"os"
//
//	tb "gopkg.in/tucnak/telebot.v2"
//)
//
//func main() {
//	var (
//		port      = os.Getenv("PORT")
//		publicURL = os.Getenv("PUBLIC_URL") // you must add it to your config vars
//		token     = os.Getenv("TOKEN")      // you must add it to your config vars
//	)
//
//	webhook := &tb.Webhook{
//		Listen:   ":" + port,
//		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
//	}
//
//	pref := tb.Settings{
//		Token:  token,
//		Poller: webhook,
//	}
//
//	b, err := tb.NewBot(pref)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	b.Start()
//
//	b.Handle("/hello", func(m *tb.Message) {
//		_, err := b.Send(m.Sender, "Hi!")
//		if err != nil {
//			log.Fatal(err)
//		}
//	})
//}

package main

import (
	"log"

	"github.com/Syfaro/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("MyAwesomeBotToken")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}