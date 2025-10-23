package main

import (
	"habit-check-bot/internal/app"
	"log"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}
	
	app.Run(logger)
}
