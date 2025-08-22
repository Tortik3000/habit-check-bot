package time_events

import "go.uber.org/zap"

type TimeEvents struct {
	logger *zap.Logger
	db     Storage
}

func New(
	logger *zap.Logger,
	db Storage,
) *TimeEvents {
	return &TimeEvents{
		logger: logger,
		db:     db,
	}
}
