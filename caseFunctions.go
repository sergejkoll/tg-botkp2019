package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Syfaro/telegram-bot-api"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"tg-botkp2019/models"
)

var userArray []models.User
var taskArray []models.Task

var (
	tokensMap = make(map[int64]models.Tokens)
)
//
// РЕГИСТРАЦИЯ
//
func getUserIdAndAddInArrayCase(bot *tgbotapi.BotAPI, id int64) {
	currentUser := models.User{
		Id: int(id),
	}
	userArray = append(userArray, currentUser)
	msg := tgbotapi.NewMessage(id, "Введите почту")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getEmailCase(bot *tgbotapi.BotAPI, id int64, email string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Email = email
	// Отправка сообщения
	msg := tgbotapi.NewMessage(id, "Введите логин")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getLoginCase(bot *tgbotapi.BotAPI, id int64, login string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Login = login
	// Отправка сообщения
	msg := tgbotapi.NewMessage(id, "Введите пароль")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getPasswordAndRegister(bot *tgbotapi.BotAPI, id int64, password string) (status bool) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Password = password
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}
	// Отправка запроса
	resp, err := http.Post("http://jtdi.ru/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	var access *http.Cookie
	var refresh *http.Cookie
	for _, cookie := range resp.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		if cookie.Name == "access_token" {
			access = cookie
		}
		if cookie.Name == "refresh_token" {
			refresh = cookie
		}
	}

	tokens := models.Tokens{
		Access:access,
		Refresh:refresh,
	}

	tokensMap[id] = tokens

	if resp.StatusCode == 200 {
		// Отпарвка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(id, "Выберите объект с которым хотите продолжить работу")
		msg.ReplyMarkup = mainMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(id, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

//
// ВХОД
//
func getPasswordAndLogin(bot *tgbotapi.BotAPI, id int64, password string) (status bool) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(id) {
			index = i
			break
		}
	}
	userArray[index].Password = password
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}
	// Отправка запроса
	resp, err := http.Post("http://jtdi.ru/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	var access *http.Cookie
	var refresh *http.Cookie
	for _, cookie := range resp.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		if cookie.Name == "access_token" {
			access = cookie
		}
		if cookie.Name == "refresh_token" {
			refresh = cookie
		}
	}

	tokens := models.Tokens{
		Access:access,
		Refresh:refresh,
	}

	tokensMap[id] = tokens

	if resp.StatusCode == 200 {
		// Отпарвка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(id, "Выберите объект с которым хотите продолжить работу")
		msg.ReplyMarkup = mainMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(id, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func getUserIdForLogin(bot *tgbotapi.BotAPI, id int64) {
	currentUser := models.User{
		Id: int(id),
	}
	userArray = append(userArray, currentUser)
	msg := tgbotapi.NewMessage(id, "Введите логин")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

//
// ЗАДАЧИ
//
func getUserIdForTask(bot *tgbotapi.BotAPI, id int64) {
	currentTask := models.Task{
		CreatorId: int(id),
	}
	taskArray = append(taskArray, currentTask)
	msg := tgbotapi.NewMessage(id, "Введите название задачи")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskTitle(bot *tgbotapi.BotAPI, userId int64, taskTitle string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(userId) {
			index = i
			break
		}
	}
	taskArray[index].Title = taskTitle
	// Отправка сообщения
	msg := tgbotapi.NewMessage(userId, "Введите крайний срок выполнения задачи в формате dd-mm-yyyy hh:mm:ss")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskDeadline(bot *tgbotapi.BotAPI, chatId int64, taskDeadline string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	//yy, mm, dd := time.Now().Date()
	//year := strconv.Itoa(yy)
	//day := strconv.Itoa(dd)
	//timeStr := year + " " + mm.String() + " " + day
	//duration, _ := time.Parse("2006 January 02 15:04", timeStr + " " + taskDuration)
	t, err := time.Parse("02-01-2006 15:04:05", taskDeadline)
	if err != nil {
		log.Fatal(err)
	}
	taskArray[index].Deadline = t.Unix()
	// Отправка сообщение о времени выполнении
	msg := tgbotapi.NewMessage(chatId, "Введите предположительное время выполнение в формате hh:mm")
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskDuration(bot *tgbotapi.BotAPI, chatId int64, taskDuration string) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	// Формирование времени
	t, _ := time.Parse("15:04", taskDuration)
	var dur int64
	dur = int64((t.Hour() * 60 * 60) + (t.Minute() * 60))
	taskArray[index].Duration = dur
	// Отправка сообщения
	msg := tgbotapi.NewMessage(chatId, "Введите приоретет задачи (от 1 до 3)")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getTaskPriority(bot *tgbotapi.BotAPI, chatId int64, taskPriority string) (status bool) {
	// Получение задачи из массива
	index := 0
	for i, task := range taskArray {
		if task.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	priority, _ := strconv.Atoi(taskPriority)
	taskArray[index].Priority = priority
	// Формирование тела запроса
	body, err := json.Marshal(taskArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/task/create"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	// Устанвока куки
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}


	if resp.StatusCode == 200 {
		// Отпарвка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(chatId, "Задача создана")
		msg.ReplyMarkup = createTaskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}
