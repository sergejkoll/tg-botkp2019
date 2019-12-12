package main

import (
	"github.com/Syfaro/telegram-bot-api"
)

var startKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Вход", "login"),
		tgbotapi.NewInlineKeyboardButtonData("Регистрация", "signup"),
	),
)

var mainMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Пользователь", "user"),
		tgbotapi.NewInlineKeyboardButtonData("Задачи", "task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Группы", "group"),
		tgbotapi.NewInlineKeyboardButtonData("Интервалы", "scope"),
	),
)

var taskMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать задачу", "create_task"),
		tgbotapi.NewInlineKeyboardButtonData("Показать задачи", "get_tasks"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновить задачу", "update_task"),
		tgbotapi.NewInlineKeyboardButtonData("Чеклисты", "checklist"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Треки", "track"),
		tgbotapi.NewInlineKeyboardButtonData("Меню", "menu"),
	),
)

var scopeMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать", "create_scope"),
		tgbotapi.NewInlineKeyboardButtonData("Удалить", "get_scope"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Редактировать", "update_scope"),
		tgbotapi.NewInlineKeyboardButtonData("Показать все", "get_allScopes"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("УМНЫЙ АЛГОРИТМ", "iftellect"),
	),
)
var createTaskMenuKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Обновить задачу", "update_task"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить чеклист", "checklist"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить задачу в группу", "add_task_in_group"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Меню", "menu"),
	),
)


