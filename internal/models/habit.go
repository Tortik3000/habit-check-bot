package models

import (
	"github.com/go-telegram/bot/models"
)

type Habit struct {
	Name   string
	ChatId int64
	Mark   bool
}

func CreateInlineKeyboard(habit *Habit) []models.InlineKeyboardButton {
	name := getHabitName(habit)
	return []models.InlineKeyboardButton{
		{Text: name, CallbackData: habit.Name},
	}

}

func getHabitName(habit *Habit) string {
	if habit.Mark {
		return habit.Name + " ✅"
	}
	return habit.Name
}
