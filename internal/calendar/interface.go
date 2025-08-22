package calendar

import (
	"context"
)

type Storage interface {
	GetMarkDay(
		ctx context.Context,
		chatId int64,
		date string,
	) (bool, error)
}
