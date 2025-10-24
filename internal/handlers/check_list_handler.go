package handlers

import (
	"context"
	m "habit-check-bot/internal/models"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) CheckListHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.logger.Info("started GetAccountsHabits")
	chatID := update.Message.Chat.ID
	habits, err := h.db.GetAccountsHabits(ctx, chatID)
	if err != nil {
		h.logger.Error(err.Error())
		return
	}

	inlineKeyboard := make([][]models.InlineKeyboardButton, 0, len(habits))
	for _, habit := range habits {
		inlineKeyboard = append(inlineKeyboard, m.CreateInlineKeyboard(&habit))
	}

	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboard,
	}

	message, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Your check list",
		ReplyMarkup: kb,
	})

	if err != nil {
		h.logger.Error(err.Error())
		return
	}
	h.logger.Info("send msg", zap.String("msg", message.Text))
}
