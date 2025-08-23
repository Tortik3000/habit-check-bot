package calendar

import (
	"fmt"
	"time"

	"github.com/go-telegram/bot/models"
)

func BuildCalendarKeyboard(month, year int, marks []bool) *models.InlineKeyboardMarkup {
	keyboard := [][]models.InlineKeyboardButton{buildCalendarHeader()}

	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	weekday := getFirstWeekday(firstDay)

	rows := buildCalendarRows(month, year, weekday, marks)

	keyboard = append(keyboard, rows...)
	keyboard = append(keyboard, buildMonthNavigation(month, year))
	return &models.InlineKeyboardMarkup{InlineKeyboard: keyboard}
}

func buildCalendarHeader() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{Text: "Mo", CallbackData: "IGNORE"},
		{Text: "Tu", CallbackData: "IGNORE"},
		{Text: "We", CallbackData: "IGNORE"},
		{Text: "Th", CallbackData: "IGNORE"},
		{Text: "Fr", CallbackData: "IGNORE"},
		{Text: "Sa", CallbackData: "IGNORE"},
		{Text: "Su", CallbackData: "IGNORE"},
	}
}

func getFirstWeekday(firstDay time.Time) int {
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	return weekday
}

func buildCalendarRows(month, year int, weekday int, marks []bool) [][]models.InlineKeyboardButton {
	var rows [][]models.InlineKeyboardButton
	row := make([]models.InlineKeyboardButton, 0, 7)
	for i := 1; i < weekday; i++ {
		row = append(row, models.InlineKeyboardButton{Text: "-", CallbackData: "IGNORE"})
	}

	for day, mark := range marks {
		row = append(row, buildDayButton(day+1, month, year, mark))
		if len(row) == 7 {
			rows = append(rows, row)
			row = []models.InlineKeyboardButton{}
		}
	}
	if len(row) > 0 {
		for len(row) < 7 {
			row = append(row, models.InlineKeyboardButton{Text: "-", CallbackData: "IGNORE"})
		}
		rows = append(rows, row)
	}
	return rows
}

func buildDayButton(day, month, year int, mark bool) models.InlineKeyboardButton {
	var btnText string
	if mark {
		btnText = fmt.Sprintf("%d 🟢", day)
	} else {
		btnText = fmt.Sprintf("%d", day)
	}
	return models.InlineKeyboardButton{
		Text:         btnText,
		CallbackData: fmt.Sprintf("DAY_%d_%d_%d", day, month, year),
	}
}

func buildMonthNavigation(month, year int) []models.InlineKeyboardButton {
	prevMonth, prevYear := month-1, year
	if prevMonth < 1 {
		prevMonth = 12
		prevYear--
	}
	nextMonth, nextYear := month+1, year
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}
	return []models.InlineKeyboardButton{
		{
			Text:         "⬅️",
			CallbackData: fmt.Sprintf("CAL_%d_%d", prevMonth, prevYear),
		},
		{
			Text:         fmt.Sprintf("%d-%02d", year, month),
			CallbackData: "IGNORE",
		},
		{
			Text:         "➡️",
			CallbackData: fmt.Sprintf("CAL_%d_%d", nextMonth, nextYear),
		},
	}
}
