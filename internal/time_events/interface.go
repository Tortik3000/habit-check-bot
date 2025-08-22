package time_events

import (
	"context"
	"habit-check-bot/internal/models"
)

type Storage interface {
	GetAccountsHabits(
		ctx context.Context,
		chatId int64,
	) ([]models.Habit, error)

	GetAllChatIDs(
		ctx context.Context,
	) ([]int64, error)

	MarkCalendarDay(
		ctx context.Context,
		chatId int64,
		date string,
		mark bool,
	) error
}
