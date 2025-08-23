package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"habit-check-bot/internal/calendar"
)

func (h *Handler) CalendarHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var month, year int
	var chatId int64
	if update.CallbackQuery != nil {
		fmt.Sscanf(update.CallbackQuery.Data, "CAL_%d_%d", &month, &year)
		chatId = update.CallbackQuery.Message.Message.Chat.ID
	} else {
		now := time.Now()
		month = int(now.Month())
		year = now.Year()
		chatId = update.Message.Chat.ID
	}

	marks := getDaysWithMark(month, year, chatId, h.db)
	kb := calendar.BuildCalendarKeyboard(month, year, marks)

	text := fmt.Sprintf("Calendar: %d-%02d", year, month)
	if update.CallbackQuery != nil {
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      chatId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        text,
			ReplyMarkup: kb,
		})
		if err != nil {
			h.logger.Error(err.Error())
		}
	} else {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      update.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: kb,
		})
		if err != nil {
			h.logger.Error(err.Error())
		}
	}
}

func getDaysWithMark(month, year int, chatId int64, db Storage) []bool {
	daysInMonth := time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
	marks := make([]bool, daysInMonth)

	for day := range marks {
		date := fmt.Sprintf("%d-%02d-%02d", year, month, day+1)
		mark, err := db.GetMarkDay(context.Background(), chatId, date)
		marks[day] = (err == nil && mark)
	}
	return marks
}
