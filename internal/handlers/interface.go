package handlers

import (
	"context"
	"habit-check-bot/internal/models"
)

type Storage interface {
	SaveHabit(
		ctx context.Context,
		name string,
		chatId int64,
	) error

	DeleteHabit(
		ctx context.Context,
		name string,
		chatId int64,
	) error

	MarkHabit(
		ctx context.Context,
		name string,
		chatId int64,
	) error

	GetHabitsDates(
		ctx context.Context,
		name string,
		chatId int64,
	) ([]string, error)

	GetAccountsHabits(
		ctx context.Context,
		chatId int64,
	) ([]models.Habit, error)

	GetMarkDay(
		ctx context.Context,
		chatId int64,
		date string,
	) (bool, error)

	MarkCalendarDay(
		ctx context.Context,
		chatId int64,
		date string,
		mark bool,
	) error
}
