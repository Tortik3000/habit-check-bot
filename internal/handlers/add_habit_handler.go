package handlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	m "habit-check-bot/internal/models"
)

func (s *Service) AddHabitHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a habit name: /add {habit_name}",
		})
		return
	}
	habitName := strings.TrimSpace(parts[1])
	habit := m.Habit{
		Name:   habitName,
		ChatId: update.Message.Chat.ID,
	}

	err := s.db.SaveHabit(ctx, &habit)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to add habit.",
		})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Habit added: " + habitName,
	})
}
