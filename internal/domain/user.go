package domain

import "time"

// User представляет пользователя Telegram
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	ChatID    int64     `json:"chat_id" db:"chat_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// FullName возвращает полное имя пользователя
func (u User) FullName() string {
	if u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	return u.FirstName
}

// DisplayName возвращает отображаемое имя с username
func (u User) DisplayName() string {
	name := u.FullName()
	if u.Username != "" {
		name += " (@" + u.Username + ")"
	}
	return name
}
