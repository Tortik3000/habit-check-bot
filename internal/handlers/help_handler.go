package handlers

import (
	"context"
	"habit-check-bot/internal/messages"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (h *Handler) HelpHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   messages.Help,
	})
	if err != nil {
		h.logger.Error("start handler", zap.Error(err))
		return
	}
}
