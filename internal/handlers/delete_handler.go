package handlers

import (
	"context"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (h *Handler) DeleteHabitHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	chatId := update.Message.Chat.ID
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Please provide a habit name: /del {habit_name}",
		})
		return
	}
	habitName := strings.TrimSpace(parts[1])

	err := h.db.DeleteHabit(ctx, habitName, chatId)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   "Failed to delete habit.",
		})
		return
	}

	habits, err := h.db.GetAccountsHabits(ctx, chatId)
	if err != nil {
		h.logger.Error(err.Error())
		return
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Habit deleted: " + habitName,
	})
}
