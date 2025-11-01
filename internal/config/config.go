package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения
type Config struct {
	BotToken string
	DBPath   string
	Debug    bool
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN environment variable is required")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "bot.db"
	}

	debug := false
	if debugStr := os.Getenv("DEBUG"); debugStr != "" {
		var err error
		debug, err = strconv.ParseBool(debugStr)
		if err != nil {
			return nil, fmt.Errorf("invalid DEBUG value: %w", err)
		}
	}

	return &Config{
		BotToken: botToken,
		DBPath:   dbPath,
		Debug:    debug,
	}, nil
}
