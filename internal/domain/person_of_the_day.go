package domain

import "time"

// PersonOfTheDay представляет запись о человеке дня
type PersonOfTheDay struct {
	ID        int       `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	Date      time.Time `json:"date" db:"date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// UserStats представляет статистику пользователя
type UserStats struct {
	User  User `json:"user"`
	Count int  `json:"count" db:"count"`
}
