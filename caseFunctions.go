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

	"github.com/tg-botkp2019/models"
)

var userArray []models.User
var taskArray []models.Task
var groupArray []models.Group

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
		_, _ = bot.Send(tgbotapi.NewMessage(id, "Бот сломался, перезапустите его"))
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
		_, _ = bot.Send(tgbotapi.NewMessage(id, "Бот сломался, перезапустите его"))
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
	userArray[index].Fullname = login
	userArray[index].Login = login
	// Отправка сообщения
	msg := tgbotapi.NewMessage(id, "Введите пароль")
	_, err := bot.Send(msg)
	if err != nil {
		_, _ = bot.Send(tgbotapi.NewMessage(id, "Бот сломался, перезапустите его"))
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
		Access: access,
		Refresh: refresh,
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
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
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

	var result models.JsonTask
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	taskArray[index].Id = result.Task.Id


	if resp.StatusCode == 200 {
		// Отправка сообщения с новой клавиатурой
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

func GetTasks(bot *tgbotapi.BotAPI, chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Список всех Ваших задач:")
	bot.Send(msg)

	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/tasks"
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonTasks
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)

	flag := false
	for _, value := range result.Tasks {
		flag = true
		msg = tgbotapi.NewMessage(chatId, "Задача № " + strconv.Itoa(value.Id) + ". " + value.Title + "\n" + value.Description)
		bot.Send(msg)
	}
	if flag {
		msg = tgbotapi.NewMessage(chatId, "Вы можете более подробно рассмотреть каждую из задач")
		bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Для этого необходимо отправить номер желаемой задачи")
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatId, "Похоже у Вас еще нет активных задач")
		msg.ReplyMarkup = taskMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//func AskNewTitle(bot *tgbotapi.BotAPI, chatId string) {
//
//	msg := tgbotapi.NewMessage(id, "Введите новое название задачи")
//	_, err := bot.Send(msg)
//	if err != nil {
//		log.Fatal(err)
//	}
//}

//
// USER
//
func GetUser(bot *tgbotapi.BotAPI, chatId int64) (status bool) {
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("GET", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}


	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result models.JsonUserBody
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)


	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Email: " + result.User.Email + "\n Login: " +
			result.User.Login + "\n" + "Fullname: " + result.User.Fullname)
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return true
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func DeleteUser(bot *tgbotapi.BotAPI, chatId int64) (status bool){
	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("DELETE", url, nil)
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Пользователь удален. Войдите или зарегестрируйтесь")
		msg.ReplyMarkup = startKeyboard
		_, _ = bot.Send(msg)
		return true
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return false
	}
}

func getNewUserEmail(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новую почту")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updateEmail(bot *tgbotapi.BotAPI,chatId int64, email string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Email = email
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Почта обновлена")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNewUserLogin(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новый логин")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updateLogin(bot *tgbotapi.BotAPI,chatId int64, login string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Login = login
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Логин обновлен")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNewUserFullname(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новое имя")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updateFullname(bot *tgbotapi.BotAPI,chatId int64, fullname string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Fullname = fullname
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Имя обновлено")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getNewUserPass(bot *tgbotapi.BotAPI,chatId int64) {
	msg := tgbotapi.NewMessage(chatId, "Введите новый пароль")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func updatePass(bot *tgbotapi.BotAPI,chatId int64, pass string) {
	// Получение пользователя из массива
	index := 0
	for i, usr := range userArray {
		if usr.Id == int(chatId) {
			index = i
			break
		}
	}
	userArray[index].Password = pass
	// Формирование запроса
	body, err := json.Marshal(userArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/user/" + idStr
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	req.AddCookie(tokensMap[chatId].Access)
	req.AddCookie(tokensMap[chatId].Refresh)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == 200 {
		msg := tgbotapi.NewMessage(chatId, "Пароль обновлен")
		msg.ReplyMarkup = userMenuKeyboard
		_, err := bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = userMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}

//
// GROUP
//
func getIdAndGroupTitle(bot *tgbotapi.BotAPI, chatId int64) {
	currentGroup := models.Group{
		CreatorId: int(chatId),
	}
	groupArray = append(groupArray, currentGroup)
	msg := tgbotapi.NewMessage(chatId, "Введите название группы")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getGroupTitle(bot *tgbotapi.BotAPI,chatId int64, title string) {
	// Получение группы из массива
	index := 0
	for i, group := range groupArray {
		if group.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	groupArray[index].Title = title
	// Отправка сообщения
	msg := tgbotapi.NewMessage(chatId, "Введите описание")
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getGroupDescriptionAndCreate(bot *tgbotapi.BotAPI,chatId int64, description string) {
	// Получение группы из массива
	index := 0
	for i, group := range groupArray {
		if group.CreatorId == int(chatId) {
			index = i
			break
		}
	}
	groupArray[index].Description = description
	// Формирование тела запроса
	body, err := json.Marshal(groupArray[index])
	if err != nil {
		log.Fatal(err)
	}

	// Формирование запроса
	idStr := strconv.Itoa(int(chatId))
	url := "http://jtdi.ru/" + idStr + "/group/create"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	// Установка куки
	if tokensMap[chatId].Access == nil || tokensMap[chatId].Refresh == nil {
		msg := tgbotapi.NewMessage(chatId, "Авторизуйтесь!")
		msg.ReplyMarkup = startKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
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

	var result models.JsonGroup
	json.NewDecoder(bytes.NewBuffer(body)).Decode(&result)
	groupArray[index].Id = result.Group.Id

	if resp.StatusCode == 200 {
		// Отправка сообщения с новой клавиатурой
		msg := tgbotapi.NewMessage(chatId, "Группа создана\nВыберите действие для группы")
		msg.ReplyMarkup = groupMenuKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		// В случе ошибки предложить еще раз
		msg := tgbotapi.NewMessage(chatId, string(body))
		_, _ = bot.Send(msg)
		msg = tgbotapi.NewMessage(chatId, "Эхх не получилось... Попробуйте еще раз!!!")
		msg.ReplyMarkup = groupCreateKeyboard
		_, err = bot.Send(msg)
		if err != nil {
			log.Fatal(err)
		}
	}
}
