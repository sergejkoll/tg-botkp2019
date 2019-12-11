package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/Syfaro/telegram-bot-api"
)

func main() {
	proxyUrl, err := url.Parse("http://51.158.123.35:8811")
	if err != nil {
		log.Println(err)
	}
	http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	bot, err := tgbotapi.NewBotAPI("1041039490:AAGXBA0Kno3_lpYlIruQ_HzgD18kW9vCYzI")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	caseState := make(map[int64]int)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет. Я телеграм-бот. Войдите или зарегистрируйтесь")
				msg.ReplyMarkup = startKeyboard
				caseState[update.Message.Chat.ID] = START
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			case "reset":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что ж, начнём по-новой) Войдите или зарегистрируйтесь")
				msg.ReplyMarkup = startKeyboard
				caseState[update.Message.Chat.ID] = START
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			case "":
				break
			}
			// Получение почты и отправка сообщения о логине
			if caseState[update.Message.Chat.ID] == REGISTER_ENTER_EMAIL {
				getEmailCase(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = REGISTER_ENTER_LOGIN
			//	Получение логина и отправка сообщения о пароле
			} else if caseState[update.Message.Chat.ID] == REGISTER_ENTER_LOGIN {
				getLoginCase(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = REGISTER_ENTER_PASS
			//	Получение пароля и формирования запроса (в случае ошибки возврат в стартовое меню)
			} else if caseState[update.Message.Chat.ID] == REGISTER_ENTER_PASS {
				status := getPasswordAndRegister(bot, update.Message.Chat.ID, update.Message.Text)
				if !status {
					caseState[update.Message.Chat.ID] = START
				}
			//	Получение логина для входа и отправка сообщения о пароле
			} else if caseState[update.Message.Chat.ID] == SIGNIN_ENTER_LOGIN {
				getLoginCase(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = SIGNIN_ENTER_PASSWORD
			//	Получение пароля для входа и формирование запроса (в случе ошибки возврат в стартовое меню)
			} else if caseState[update.Message.Chat.ID] == SIGNIN_ENTER_PASSWORD {
				status := getPasswordAndLogin(bot, update.Message.Chat.ID, update.Message.Text)
				if !status {
					caseState[update.Message.Chat.ID] = START
				}
			//	Получение заголовка задачи и отправка сообщения о деадлайне
			} else if caseState[update.Message.Chat.ID] == TASK_SEND_TITLE {
				getTaskTitle(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = TASK_SEND_DEADLINE
			//	Получение деадлайна задачи и отправка сообщения о продолжительности
			} else if caseState[update.Message.Chat.ID] == TASK_SEND_DEADLINE {
				getTaskDeadline(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = TASK_SEND_DURATION
			//	Получение продолжительности и отапрвка сообщения о приоретете
			} else if caseState[update.Message.Chat.ID] == TASK_SEND_DURATION {
				getTaskDuration(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = TASK_SEND_PRIORITY
			//	Получение приоретета и формирование запроса
			} else if caseState[update.Message.Chat.ID] == TASK_SEND_PRIORITY {
				status := getTaskPriority(bot, update.Message.Chat.ID, update.Message.Text)
				if !status {
					caseState[update.Message.Chat.ID] = MENU
				}
			}
		}

		if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case "signup":
				getUserIdAndAddInArrayCase(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = REGISTER_ENTER_EMAIL
			case "login":
				getUserIdForLogin(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = SIGNIN_ENTER_LOGIN
			case "task":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выбирите действие для задачи")
				msg.ReplyMarkup = taskMenuKeyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			case "create_task":
				getUserIdForTask(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = TASK_SEND_TITLE
			case "menu":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите объект с которым хотите продолжить работу")
				msg.ReplyMarkup = mainMenuKeyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
