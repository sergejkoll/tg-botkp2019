package main

import (
	"bytes"
	"log"
	"net/http"
	"net/url"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	requestMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет. Я телеграм-бот. Войдите или зарегистрируйтесь")
	loginButton := tgbotapi.NewKeyboardButton("/signin")
	registerButton := tgbotapi.NewKeyboardButton("/signup")
	keyboard := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{loginButton, registerButton})
	requestMessage.ReplyMarkup = keyboard
	_, err := bot.Send(requestMessage)
	if err != nil {
		log.Fatal(err)
	}
}

func login(bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	log.Println(update.Message.Chat.ID)

	requestMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите логин")
	_, err := bot.Send(requestMessage)
	if err != nil {
		log.Fatal(err)
	}
	var login string
	for upd := range updates {
		if upd.Message != nil {
			login = upd.Message.Text
			log.Println("This is user login: " + login)
			break
		}
	}

	log.Println(update.Message.Chat.ID)

	requestMessage.Text = "Введите пароль"
	_, err = bot.Send(requestMessage)
	if err != nil {
		log.Fatal(err)
	}
	var pass string
	for upd := range updates {
		if upd.Message != nil {
			pass = upd.Message.Text
			log.Println("This is user pass" + pass)
			break
		}
	}


	json := `{"login":"` + login + `", "password":"` + pass + `"}`
	log.Println(json)
	body := bytes.NewBuffer([]byte(json))

	resp, err := http.Post("http://127.0.0.1:8080/login", "application/json", body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)
}


func GoBot() {
	proxyUrl, err := url.Parse("http://51.158.123.35:8811")
	if err != nil {
		log.Println(err)
	}
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	// используя токен создаем новый инстанс бота
	bot, err := tgbotapi.NewBotAPI("1041039490:AAGXBA0Kno3_lpYlIruQ_HzgD18kW9vCYzI")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// u - структура с конфигом для получения апдейтов
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// используя конфиг u создаем канал в который будут прилетать новые сообщения
	updates, err := bot.GetUpdatesChan(u)

	// в канал updates прилетают структуры типа Update
	// вычитываем их и обрабатываем
	for update := range updates {
		if update.Message == nil {
			continue
		}
		// логируем от кого какое сообщение пришло
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		// свитч на обработку комманд
		// комманда - сообщение, начинающееся с "/"
		switch update.Message.Command() {
		case "start":
			start(bot, update)
		case "signin":
			login(bot, update, updates)
		case "signup":
			log.Println(update.Message.Chat.ID)
		}
	}
}

func main() {
	GoBot()
}