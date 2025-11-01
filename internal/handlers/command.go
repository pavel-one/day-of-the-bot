package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/pavel-one/day-of-the-bot/internal/repository"
	"github.com/pavel-one/day-of-the-bot/internal/templates"
	"gopkg.in/telebot.v3"
)

// CommandHandler обрабатывает команды бота
type CommandHandler struct {
	api                *telebot.Bot
	userRepo           repository.UserRepository
	personOfTheDayRepo repository.PersonOfTheDayRepository
	messageService     *templates.MessageService
	rng                *rand.Rand
}

// NewCommandHandler создает новый обработчик команд
func NewCommandHandler(
	api *telebot.Bot,
	userRepo repository.UserRepository,
	personOfTheDayRepo repository.PersonOfTheDayRepository,
	messageService *templates.MessageService,
	rng *rand.Rand,
) *CommandHandler {
	return &CommandHandler{
		api:                api,
		userRepo:           userRepo,
		personOfTheDayRepo: personOfTheDayRepo,
		messageService:     messageService,
		rng:                rng,
	}
}

// RegisterHandlers регистрирует обработчики команд
func (h *CommandHandler) RegisterHandlers(bot *telebot.Bot) {
	bot.Handle("/start", h.handleStart)
	bot.Handle("/help", h.handleStart)
	bot.Handle("/pidor", h.handlePersonOfTheDay)
	bot.Handle("/pidorstats", h.handleStats)
	bot.Handle("/pidorinfo", h.handleInfo)
}

func (h *CommandHandler) handleStart(c telebot.Context) error {
	log.Printf("Команда /start вызвана в чате %d пользователем %d", c.Chat().ID, c.Sender().ID)
	SafeSendMessage(c, h.messageService.HelpText())
	return nil
}

func (h *CommandHandler) handlePersonOfTheDay(c telebot.Context) error {
	log.Printf("Команда /pidor вызвана в чате %d пользователем %d", c.Chat().ID, c.Sender().ID)
	// Проверяем, выбран ли пидор дня на сегодня
	todayPerson, err := h.personOfTheDayRepo.GetByDate(c.Chat().ID, time.Now())
	if err != nil {
		log.Printf("Ошибка при проверке пидора дня: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при проверке пидора дня"))
		return nil
	}

	if todayPerson != nil {
		SafeSendMessage(c, h.messageService.PersonAlreadySelected(*todayPerson))
		return nil
	}

	// Получаем список активных пользователей
	users, err := h.userRepo.GetByChatID(c.Chat().ID)
	if err != nil {
		log.Printf("Ошибка при получении списка участников: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при получении списка участников"))
		return nil
	}

	if len(users) == 0 {
		SafeSendMessage(c, h.messageService.NoActiveUsers())
		return nil
	}

	// Выбираем случайного пользователя
	selectedUser := users[h.rng.Intn(len(users))]

	// Сохраняем результат
	err = h.personOfTheDayRepo.Set(selectedUser.ID, c.Chat().ID, time.Now())
	if err != nil {
		log.Printf("Ошибка при сохранении результата: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при сохранении результата"))
		return nil
	}

	SafeSendMessage(c, h.messageService.PersonSelected(selectedUser))
	return nil
}

func (h *CommandHandler) handleStats(c telebot.Context) error {
	stats, err := h.personOfTheDayRepo.GetUserStats(c.Chat().ID)
	if err != nil {
		log.Printf("Ошибка при получении статистики: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при получении статистики"))
		return nil
	}

	if len(stats) == 0 {
		SafeSendMessage(c, "Статистика пока пуста.")
		return nil
	}

	SafeSendMessage(c, h.messageService.BuildStatsMessage(stats))
	return nil
}

func (h *CommandHandler) handleInfo(c telebot.Context) error {
	// Логируем информацию о сообщении для debug
	log.Printf("Команда /pidorinfo вызвана в чате %d пользователем %d", c.Chat().ID, c.Sender().ID)

	// Получаем статистику чата
	stats, err := h.personOfTheDayRepo.GetUserStats(c.Chat().ID)
	if err != nil {
		log.Printf("Ошибка при получении статистики: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при получении информации"))
		return nil
	}

	// Получаем количество активных пользователей
	users, err := h.userRepo.GetByChatID(c.Chat().ID)
	if err != nil {
		log.Printf("Ошибка при получении списка пользователей: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при получении информации о пользователях"))
		return nil
	}

	// Проверяем, выбран ли пидор на сегодня
	todayPerson, err := h.personOfTheDayRepo.GetByDate(c.Chat().ID, time.Now())
	if err != nil {
		log.Printf("Ошибка при проверке пидора дня: %v", err)
		SafeSendMessage(c, h.messageService.ErrorOccurred("при проверке пидора дня"))
		return nil
	}

	// Формируем информационное сообщение
	infoMsg := fmt.Sprintf("📊 Информация о чате:\n\n")
	infoMsg += fmt.Sprintf("👥 Активных пользователей: %d\n", len(users))
	infoMsg += fmt.Sprintf("🏆 Записей в статистике: %d\n", len(stats))

	if todayPerson != nil {
		infoMsg += fmt.Sprintf("🎯 Пидор дня сегодня: %s", todayPerson.FullName())
	} else {
		infoMsg += "🎯 Пидор дня сегодня еще не выбран"
	}

	SafeSendMessage(c, infoMsg)
	return nil
}
