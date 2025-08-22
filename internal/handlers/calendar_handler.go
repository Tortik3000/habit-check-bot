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
	if update.CallbackQuery != nil && update.CallbackQuery.Data != "" {
		fmt.Sscanf(update.CallbackQuery.Data, "CAL_%d_%d", &month, &year)
	} else {
		now := time.Now()
		month = int(now.Month())
		year = now.Year()
	}
	var chatId int64
	if update.Message != nil {
		chatId = update.Message.Chat.ID
	}

	kb := calendar.BuildCalendarKeyboard(month, year, chatId, h.db)

	text := fmt.Sprintf("Calendar: %d-%02d", year, month)
	if update.CallbackQuery != nil {
		_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
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
