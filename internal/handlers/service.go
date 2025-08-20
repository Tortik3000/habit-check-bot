package handlers

import (
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	db     Storage
}

func New(
	logger *zap.Logger,
	db Storage,
) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}
