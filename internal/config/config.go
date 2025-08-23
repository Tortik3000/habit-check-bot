package config

import (
	"fmt"
	"net"
	"os"
)

type (
	Config struct {
		PG
		TG
	}

	PG struct {
		URL      string
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		DB       string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
	}

	TG struct {
		Token string `env:"TELEGRAM_BOT_TOKEN"`
	}
)

func New() *Config {
	cfg := &Config{}

	cfg.TG.Token = os.Getenv("TELEGRAM_BOT_TOKEN")

	cfg.PG.Host = os.Getenv("POSTGRES_HOST")
	cfg.PG.Port = os.Getenv("POSTGRES_PORT")
	cfg.PG.DB = os.Getenv("POSTGRES_DB")
	cfg.PG.User = os.Getenv("POSTGRES_USER")
	cfg.PG.Password = os.Getenv("POSTGRES_PASSWORD")

	hostPort := net.JoinHostPort(cfg.PG.Host, cfg.PG.Port)
	cfg.PG.URL = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.PG.User,
		cfg.PG.Password,
		hostPort,
		cfg.PG.DB,
	)

	return cfg
}
