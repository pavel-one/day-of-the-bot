package repository

import (
	"time"

	"github.com/pavel-one/day-of-the-bot/internal/domain"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Add(user domain.User) error
	GetByChatID(chatID int64) ([]domain.User, error)
	GetByID(userID, chatID int64) (*domain.User, error)
}

// PersonOfTheDayRepository определяет интерфейс для работы с записями человека дня
type PersonOfTheDayRepository interface {
	Set(userID, chatID int64, date time.Time) error
	GetByDate(chatID int64, date time.Time) (*domain.User, error)
	GetUserStats(chatID int64) ([]domain.UserStats, error)
}
