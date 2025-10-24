package handlers

import (
	"context"
	m "habit-check-bot/internal/models"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	callbackData := update.CallbackQuery.Data
	chatId := update.CallbackQuery.Message.Message.Chat.ID
	h.logger.Info("start GetAccountsHabits")
	habits, err := h.db.GetAccountsHabits(ctx, chatId)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	h.logger.Info("finish GetAccountsHabits", zap.Any("habits", habits))

	inlineKeyboard := make([][]models.InlineKeyboardButton, 0, len(habits))
	for i, habit := range habits {
		if habit.Name == callbackData {
			habit.Mark = !habit.Mark
			habits[i].Mark = habit.Mark
			h.logger.Info("start MarkHabit")
			err = h.db.MarkHabit(ctx, habit.Name, habit.ChatId)
			if err != nil {
				h.logger.Error(err.Error())
				return
			}
			h.logger.Info("finish MarkHabit")
		}
		inlineKeyboard = append(inlineKeyboard, m.CreateInlineKeyboard(&habit))
	}

	allDone := true
	for _, habit := range habits {
		if !habit.Mark {
			allDone = false
			break
		}
	}

	today := time.Now().Format("2006-01-02")
	h.logger.Info("start MarkCalendarDay")
	err = h.db.MarkCalendarDay(ctx, chatId, today, allDone)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}

	h.logger.Info("finish MarkCalendarDay")

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	if err != nil {
		h.logger.Error(err.Error())
	}
}
