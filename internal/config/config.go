package config

import (
	"log"
	"sync"

	"github.com/cristalhq/aconfig"
	"github.com/cristalhq/aconfig/aconfighcl"
)

type Config struct {
	TelegramBotToken string `hcl:"telegram_bot_token" env:"TELEGRAM_BOT_TOKEN" required:"true"`
	DatabaseDSN      string `hcl:"database_dsn" env:"DATABASE_DSN" default:"postgres://ed:1234567@localhost:5432/bot?sslmode=disable"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		loader := aconfig.LoaderFor(&cfg, aconfig.Config{
			EnvPrefix: "HCB",
			Files:     []string{"./config.hcl", "./config.local.hcl", "$HOME/.config/habit-check-bot/config.hcl"},
			FileDecoders: map[string]aconfig.FileDecoder{
				".hcl": aconfighcl.New(),
			},
		})

		if err := loader.Load(); err != nil {
			log.Printf("[ERROR] failed to load config: %v", err)
		}
	})

	return cfg
}
