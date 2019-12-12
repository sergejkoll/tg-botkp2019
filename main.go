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
					caseState[update.Message.Chat.ID] = TASK_MENU
				}
			// Получение одного задания из предложенных
			} else if caseState[update.Message.Chat.ID] == GOT_ALL_TASK {
				GetTaskById(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = TASK_MENU
			// Обновление задания
			} else if caseState[update.Message.Chat.ID] == UPDATE_APPROVED {
				state := AskNewTitle(bot, update.Message.Chat.ID, update.Message.Text)
				if state == 1 {
					caseState[update.Message.Chat.ID] = START
				} else if state == 2 {
					caseState[update.Message.Chat.ID] = TASK_MENU
				} else if state == 3 {
					caseState[update.Message.Chat.ID] = GOT_TITLE
				}
			} else if caseState[update.Message.Chat.ID] == GOT_TITLE {
				GetNewTaskTitle(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = ASKED_DEADLINE
			} else if caseState[update.Message.Chat.ID] == ASKED_DEADLINE {
				GetNewTaskDeadline(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = ASKED_DURATION
			} else if caseState[update.Message.Chat.ID] == ASKED_DURATION {
				GetNewTaskDuration(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = ASKED_PRIORITY
			} else if caseState[update.Message.Chat.ID] == ASKED_PRIORITY {
				status := GetNewTaskPriority(bot, update.Message.Chat.ID, update.Message.Text)
				if !status {
					caseState[update.Message.Chat.ID] = TASK_MENU
				}
			//	Получение новой почты при обновлении пользователя
			} else if caseState[update.Message.Chat.ID] == UPDATE_USER_EMAIL {
				updateEmail(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = USER_MENU
			//	Получение и обновление логина пользователя
			} else if caseState[update.Message.Chat.ID] == UPDATE_USER_LOGIN {
				updateLogin(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = USER_MENU
			//	Получение и обновление имени пользователя
			} else if caseState[update.Message.Chat.ID] == UPDATE_USER_FULLNAME {
				updateFullname(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = USER_MENU
			//	Получение и обновление пароля пользователя
			} else if caseState[update.Message.Chat.ID] == UPDATE_USER_PASS {
				updatePass(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = USER_MENU
			//	Получение заголовка группы и отправка сообщения о описании
			} else if caseState[update.Message.Chat.ID] == GROUP_SEND_TITLE {
				getGroupTitle(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = GROUP_SEND_DESCRIPTION
			} else if caseState[update.Message.Chat.ID] == GROUP_SEND_DESCRIPTION {
				getGroupDescriptionAndCreate(bot, update.Message.Chat.ID, update.Message.Text)
				caseState[update.Message.Chat.ID] = GROUP_MENU
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
			case "get_tasks":
				GetTasks(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = GOT_ALL_TASK
			case "update_task":
				UpdateTask(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = UPDATE_APPROVED
			case "scope":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выбирите действие для интервала")
				msg.ReplyMarkup = taskMenuKeyboard
			case "user":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите действие для пользователя")
				msg.ReplyMarkup = userMenuKeyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
				caseState[update.CallbackQuery.Message.Chat.ID] = USER_MENU
			case "get_user":
				GetUser(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = USER_MENU
			case "delete_user":
				status := DeleteUser(bot, update.CallbackQuery.Message.Chat.ID)
				if !status {
					caseState[update.CallbackQuery.Message.Chat.ID] = USER_MENU
				}
				caseState[update.CallbackQuery.Message.Chat.ID] = START
			case "update_user":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите что хотите обновить у пользователя")
				msg.ReplyMarkup = updateUserKeyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			case "update_email":
				getNewUserEmail(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = UPDATE_USER_EMAIL
			case "update_login":
				getNewUserLogin(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = UPDATE_USER_LOGIN
			case "update_name":
				getNewUserFullname(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = UPDATE_USER_FULLNAME
			case "update_pass":
				getNewUserPass(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = UPDATE_USER_PASS
			case "group":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите действие для группы")
				msg.ReplyMarkup = groupMenuKeyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			case "create_group":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите что именно вы хотите создать\n" +
					"(Чтобы создать любой элемент в группе, необходимо создать группу)")
				msg.ReplyMarkup = groupCreateKeyboard
				_, err = bot.Send(msg)
				if err != nil {
					log.Fatal(err)
				}
			case "create_groups":
				getIdAndGroupTitle(bot, update.CallbackQuery.Message.Chat.ID)
				caseState[update.CallbackQuery.Message.Chat.ID] = GROUP_SEND_TITLE
			}

		}
	}
}
