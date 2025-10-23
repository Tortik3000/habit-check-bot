package app

import (
	"context"
	"habit-check-bot/internal/handlers"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"habit-check-bot/db"
	"habit-check-bot/internal/config"
	repository "habit-check-bot/internal/storage/postgres"
)

func Run(logger *zap.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.New()
	dbPool, err := pgxpool.New(ctx, cfg.PG.URL)
	if err != nil {
		logger.Error("can not create pgxpool", zap.Error(err))
		return
	}

	db.SetupPostgres(dbPool, logger)
	storage := repository.New(dbPool, logger)

	service := handlers.New(logger, storage)

	opts := []bot.Option{
		bot.WithDefaultHandler(service.DefaultHandler),
		bot.WithCallbackQueryDataHandler("CAL_", bot.MatchTypePrefix, service.CalendarHandler),
		bot.WithCallbackQueryDataHandler("IGNORE", bot.MatchTypePrefix, service.EmptyHandler),
		bot.WithCallbackQueryDataHandler("DAY", bot.MatchTypePrefix, service.EmptyHandler),
		bot.WithCallbackQueryDataHandler("", bot.MatchTypePrefix, service.CallbackHandler),
	}

	b, err := bot.New(cfg.TG.Token, opts...)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start",
		bot.MatchTypeExact, service.StartHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/get",
		bot.MatchTypeExact, service.CheckListHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/add",
		bot.MatchTypePrefix, service.AddHabitHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/del",
		bot.MatchTypePrefix, service.DeleteHabitHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/cal",
		bot.MatchTypeExact, service.CalendarHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help",
		bot.MatchTypeExact, service.HelpHandler)

	b.Start(ctx)
}
