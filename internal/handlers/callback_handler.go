package handlers

import (
	"context"
	m "habit-check-bot/internal/models"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (s *Service) CallbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	// Получить актуальный список привычек
	callbackData := update.CallbackQuery.Data

	habits := s.db.GetHabits(ctx, update.CallbackQuery.Message.Message.Chat.ID)
	inlineKeyboard := make([][]models.InlineKeyboardButton, 0, len(habits))
	for _, habit := range habits {

		if habit.Name == callbackData {
			habit.Mark = true
			s.db.MarkHabit(ctx, habit.Name, habit.ChatId, true)
		}
		inlineKeyboard = append(inlineKeyboard, m.CreateInlineKeyboard(&habit))
	}

	_, err := b.EditMessageReplyMarkup(ctx, &bot.EditMessageReplyMarkupParams{
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

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "You selected the button: " + update.CallbackQuery.Data,
	})
}
