package handlers

import (
	"context"
	m "habit-check-bot/internal/models"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	callbackData := update.CallbackQuery.Data

	habits, err := h.db.GetAccountsHabits(ctx, update.CallbackQuery.Message.Message.Chat.ID)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	inlineKeyboard := make([][]models.InlineKeyboardButton, 0, len(habits))
	for _, habit := range habits {
		if habit.Name == callbackData {
			habit.Mark = !habit.Mark
			h.db.MarkHabit(ctx, habit.Name, habit.ChatId)
		}
		inlineKeyboard = append(inlineKeyboard, m.CreateInlineKeyboard(&habit))
	}

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
