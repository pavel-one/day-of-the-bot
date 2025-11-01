package templates

import (
	"fmt"
	"strings"
	"time"

	"github.com/pavel-one/day-of-the-bot/internal/domain"
)

// MessageService предоставляет методы для форматирования сообщений
type MessageService struct {
	messages *Messages
}

// NewMessageService создает новый сервис сообщений
func NewMessageService() (*MessageService, error) {
	messages, err := NewMessages()
	if err != nil {
		return nil, err
	}

	return &MessageService{
		messages: messages,
	}, nil
}

// BotGroupOnly возвращает сообщение о работе только в группах
func (ms *MessageService) BotGroupOnly() string {
	return ms.messages.BotGroupOnly.Execute(nil)
}

// UnknownCommand возвращает сообщение о неизвестной команде
func (ms *MessageService) UnknownCommand() string {
	return ms.messages.UnknownCommand.Execute(nil)
}

// ErrorOccurred возвращает сообщение об ошибке
func (ms *MessageService) ErrorOccurred(errorMsg string) string {
	return ms.messages.ErrorOccurred.Execute(TemplateData{
		"error": errorMsg,
	})
}

// HelpText возвращает текст справки
func (ms *MessageService) HelpText() string {
	return ms.messages.HelpText.Execute(nil)
}

// PersonAlreadySelected возвращает сообщение о том, что пидор дня уже выбран
func (ms *MessageService) PersonAlreadySelected(person domain.User) string {
	return ms.messages.PersonAlreadySelected.Execute(TemplateData{
		"person": person.DisplayName(),
	})
}

// PersonSelected возвращает сообщение о выборе пидора дня
func (ms *MessageService) PersonSelected(person domain.User) string {
	return ms.messages.PersonSelected.Execute(TemplateData{
		"person": person.DisplayName(),
	})
}

// NoActiveUsers возвращает сообщение об отсутствии активных пользователей
func (ms *MessageService) NoActiveUsers() string {
	return ms.messages.NoActiveUsers.Execute(nil)
}

// PersonInfo возвращает информацию о пидоре дня
func (ms *MessageService) PersonInfo(person domain.User, date time.Time) string {
	return ms.messages.PersonInfo.Execute(TemplateData{
		"person": person.DisplayName(),
		"date":   date.Format("02.01.2006"),
	})
}

// NoPersonSelectedToday возвращает сообщение о том, что сегодня пидор не выбран
func (ms *MessageService) NoPersonSelectedToday() string {
	return ms.messages.NoPersonSelectedToday.Execute(nil)
}

// StatsEmpty возвращает сообщение об отсутствии статистики
func (ms *MessageService) StatsEmpty() string {
	return ms.messages.StatsEmpty.Execute(nil)
}

// NoStatsAvailable возвращает сообщение об отсутствии статистики (алиас для совместимости)
func (ms *MessageService) NoStatsAvailable() string {
	return ms.StatsEmpty()
}

// StatsHeader возвращает отформатированную статистику
func (ms *MessageService) StatsHeader(stats []domain.UserStats) string {
	return ms.BuildStatsMessage(stats)
}

// BuildStatsMessage строит сообщение со статистикой
func (ms *MessageService) BuildStatsMessage(stats []domain.UserStats) string {
	var result strings.Builder

	// Добавляем заголовок
	result.WriteString(ms.messages.StatsHeader.Execute(nil))

	// Добавляем записи статистики
	for i, stat := range stats {
		position := GetPositionEmoji(i + 1)
		entry := ms.messages.StatsEntry.Execute(TemplateData{
			"position": position,
			"person":   stat.User.DisplayName(),
			"count":    fmt.Sprintf("%d", stat.Count),
		})
		result.WriteString(entry)
	}

	return result.String()
}
