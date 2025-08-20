package handlers

import (
	"context"
	"habit-check-bot/internal/models"
)

type Storage interface {
	SaveHabit(
		ctx context.Context,
		habit *models.Habit,
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
		mark bool,
	) error

	GetDatesHabit(
		ctx context.Context,
		name string,
		chatId int64,
	) ([]string, error)

	GetHabits(
		ctx context.Context,
		chatId int64,
	) []models.Habit
}
