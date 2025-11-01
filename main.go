package main

import (
	"log"
	"time"

	"github.com/pavel-one/day-of-the-bot/internal/bot"
	"github.com/pavel-one/day-of-the-bot/internal/config"
	"github.com/pavel-one/day-of-the-bot/internal/repository"
	"github.com/pavel-one/day-of-the-bot/internal/templates"
	"gopkg.in/telebot.v3"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Создаем настройки для telebot
	settings := telebot.Settings{
		Token:  cfg.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	// Инициализируем бота
	api, err := telebot.NewBot(settings)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	log.Printf("Авторизован как %s", api.Me.Username)

	// Инициализируем базу данных
	db, err := repository.NewDatabase(cfg.DBPath)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	personOfTheDayRepo := repository.NewPersonOfTheDayRepository(db)

	// Создаем сервис сообщений
	messageService, err := templates.NewMessageService()
	if err != nil {
		log.Fatalf("Ошибка создания сервиса сообщений: %v", err)
	}

	// Создаем и запускаем бота
	botInstance := bot.NewBot(api, userRepo, personOfTheDayRepo, messageService)
	botInstance.Start()
}
