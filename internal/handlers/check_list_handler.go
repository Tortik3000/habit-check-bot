package handlers

import (
	"context"
	m "habit-check-bot/internal/models"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"
)

func (s *Service) CheckListHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID

	habits := s.db.GetHabits(ctx, chatID)

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
		s.logger.Error(err.Error())
		return
	}
	s.logger.Info("send msg", zap.String("msg", message.Text))

	_, err = b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.Message.Chat.ID,
		MessageID: update.Message.ID,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
	s.logger.Info("delete msg", zap.String("msg", update.Message.Text))
}
