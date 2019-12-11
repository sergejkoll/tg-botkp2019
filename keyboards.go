package main

import (
	"github.com/Syfaro/telegram-bot-api"
)

var startKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("вход", "login"),
		tgbotapi.NewInlineKeyboardButtonData("регистрация", "signup"),
	),
)

var mainMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("пользователь", "user"),
		tgbotapi.NewInlineKeyboardButtonData("задачи", "task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("группы", "group"),
		tgbotapi.NewInlineKeyboardButtonData("интервалы", "scope"),
	),
)

var taskMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("создать задачу", "create_task"),
		tgbotapi.NewInlineKeyboardButtonData("показать задачи", "get_tasks"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("обновить задачу", "update_task"),
		tgbotapi.NewInlineKeyboardButtonData("чеклисты", "checklist"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("треки", "track"),
		tgbotapi.NewInlineKeyboardButtonData("меню", "menu"),
	),
)

var createTaskMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("обновить задачу", "update_task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("добавить чеклист", "checklist"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("добавить задачу в группу", "add_task_in_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("меню", "menu"),
	),
)


