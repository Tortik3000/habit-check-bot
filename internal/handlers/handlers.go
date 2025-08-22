package handlers

import (
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	db     Storage
}

func New(
	logger *zap.Logger,
	db Storage,
) *Handler {
	return &Handler{
		logger: logger,
		db:     db,
	}
}
