package handlers

import (
	"context"
	m "habit-check-bot/internal/models"
	"log"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	callbackData := update.CallbackQuery.Data
	chatId := update.CallbackQuery.Message.Message.Chat.ID
	habits, err := h.db.GetAccountsHabits(ctx, chatId)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	inlineKeyboard := make([][]models.InlineKeyboardButton, 0, len(habits))
	for i, habit := range habits {
		if habit.Name == callbackData {
			habit.Mark = !habit.Mark
			habits[i].Mark = habit.Mark
			h.db.MarkHabit(ctx, habit.Name, habit.ChatId)
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
	h.db.MarkCalendarDay(ctx, chatId, today, allDone)

	_, err = b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
		ReplyMarkup: &models.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	})

	if err != nil {
		log.Println(err)

		return
	}
}
