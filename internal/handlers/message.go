package handlers

import (
	"log"

	"github.com/pavel-one/day-of-the-bot/internal/domain"
	"github.com/pavel-one/day-of-the-bot/internal/repository"
	"github.com/pavel-one/day-of-the-bot/internal/templates"
	"gopkg.in/telebot.v3"
)

// MessageHandler обрабатывает входящие сообщения
type MessageHandler struct {
	api                *telebot.Bot
	userRepo           repository.UserRepository
	personOfTheDayRepo repository.PersonOfTheDayRepository
	messageService     *templates.MessageService
	commandHandler     *CommandHandler
}

// NewMessageHandler создает новый обработчик сообщений
func NewMessageHandler(
	api *telebot.Bot,
	userRepo repository.UserRepository,
	personOfTheDayRepo repository.PersonOfTheDayRepository,
	messageService *templates.MessageService,
	commandHandler *CommandHandler,
) *MessageHandler {
	return &MessageHandler{
		api:                api,
		userRepo:           userRepo,
		personOfTheDayRepo: personOfTheDayRepo,
		messageService:     messageService,
		commandHandler:     commandHandler,
	}
}

// RegisterHandlers регистрирует обработчики сообщений
func (h *MessageHandler) RegisterHandlers(bot *telebot.Bot) {
	// Регистрируем middleware для обработки всех сообщений
	bot.Use(h.handleMessage)

	// Регистрируем обработчики команд
	h.commandHandler.RegisterHandlers(bot)

	// Альтернативно: регистрируем обработчик для всех текстовых сообщений
	bot.Handle(telebot.OnText, h.handleTextMessage)
}

// handleMessage обрабатывает все входящие сообщения (middleware)
func (h *MessageHandler) handleMessage(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		// Логируем для отладки
		log.Printf("Middleware: получено обновление, тип чата: %s, ID чата: %d", c.Chat().Type, c.Chat().ID)

		// Пропускаем, если нет сообщения
		if c.Message() == nil {
			log.Printf("Middleware: нет сообщения, пропускаем")
			return next(c)
		}

		// Работаем только в группах
		if c.Chat().Type != telebot.ChatGroup && c.Chat().Type != telebot.ChatSuperGroup {
			log.Printf("Middleware: приватный чат, отправляем предупреждение")
			SafeSendMessage(c, h.messageService.BotGroupOnly())
			return nil // Не продолжаем обработку для приватных чатов
		}

		log.Printf("Middleware: групповой чат, обрабатываем пользователя")

		// Добавляем пользователя в базу данных
		if c.Sender() != nil {
			user := domain.User{
				ID:        c.Sender().ID,
				Username:  c.Sender().Username,
				FirstName: c.Sender().FirstName,
				LastName:  c.Sender().LastName,
				ChatID:    c.Chat().ID,
			}

			if err := h.userRepo.Add(user); err != nil {
				log.Printf("Ошибка добавления пользователя: %v", err)
			} else {
				log.Printf("Middleware: пользователь %s добавлен/обновлен", user.FirstName)
			}
		}

		// Продолжаем выполнение следующего обработчика
		log.Printf("Middleware: передаем управление следующему обработчику")
		return next(c)
	}
}

// handleTextMessage обрабатывает текстовые сообщения (альтернативный способ)
func (h *MessageHandler) handleTextMessage(c telebot.Context) error {
	log.Printf("TextHandler: получено текстовое сообщение в чате %d от пользователя %d", c.Chat().ID, c.Sender().ID)

	// Работаем только в группах
	if c.Chat().Type != telebot.ChatGroup && c.Chat().Type != telebot.ChatSuperGroup {
		log.Printf("TextHandler: приватный чат, отправляем предупреждение")
		SafeSendMessage(c, h.messageService.BotGroupOnly())
		return nil
	}

	// Добавляем пользователя в базу данных
	if c.Sender() != nil {
		user := domain.User{
			ID:        c.Sender().ID,
			Username:  c.Sender().Username,
			FirstName: c.Sender().FirstName,
			LastName:  c.Sender().LastName,
			ChatID:    c.Chat().ID,
		}

		if err := h.userRepo.Add(user); err != nil {
			log.Printf("TextHandler: Ошибка добавления пользователя: %v", err)
		} else {
			log.Printf("TextHandler: пользователь %s добавлен/обновлен", user.FirstName)
		}
	}

	return nil
}
