package main

import (
	"os"
	"testing"
	"time"

	"github.com/pavel-one/day-of-the-bot/internal/domain"
	"github.com/pavel-one/day-of-the-bot/internal/repository"
)

func TestDatabase(t *testing.T) {
	// Создаем временную базу данных для тестов
	dbPath := "test.db"
	defer func() {
		if err := os.Remove(dbPath); err != nil {
			t.Logf("Не удалось удалить тестовую БД: %v", err)
		}
	}()

	db, err := repository.NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Ошибка создания базы данных: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Logf("Ошибка закрытия БД: %v", err)
		}
	}()

	// Создаем репозитории
	userRepo := repository.NewUserRepository(db)
	personRepo := repository.NewPersonOfTheDayRepository(db)

	// Тестируем добавление пользователя
	user := domain.User{
		ID:        12345,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		ChatID:    -123456789,
	}

	err = userRepo.Add(user)
	if err != nil {
		t.Fatalf("Ошибка добавления пользователя: %v", err)
	}

	// Тестируем получение пользователей чата
	users, err := userRepo.GetByChatID(user.ChatID)
	if err != nil {
		t.Fatalf("Ошибка получения пользователей чата: %v", err)
	}

	if len(users) != 1 {
		t.Errorf("Ожидался 1 пользователь, получено %d", len(users))
	}

	if users[0].ID != user.ID {
		t.Errorf("Ожидался пользователь с ID %d, получен %d", user.ID, users[0].ID)
	}

	// Тестируем установку человека дня
	now := time.Now()
	err = personRepo.Set(user.ID, user.ChatID, now)
	if err != nil {
		t.Fatalf("Ошибка установки человека дня: %v", err)
	}

	// Тестируем получение сегодняшнего человека дня
	todayPerson, err := personRepo.GetByDate(user.ChatID, now)
	if err != nil {
		t.Fatalf("Ошибка получения человека дня: %v", err)
	}

	if todayPerson == nil {
		t.Fatal("Пидор дня должен быть установлен")
	}

	if todayPerson.ID != user.ID {
		t.Errorf("Ожидался пидор дня с ID %d, получен %d", user.ID, todayPerson.ID)
	}

	// Тестируем получение статистики
	stats, err := personRepo.GetUserStats(user.ChatID)
	if err != nil {
		t.Fatalf("Ошибка получения статистики: %v", err)
	}

	if len(stats) != 1 {
		t.Errorf("Ожидалась статистика для 1 пользователя, получено %d", len(stats))
	}

	if stats[0].Count != 1 {
		t.Errorf("Ожидалось количество 1, получено %d", stats[0].Count)
	}
}

func TestFormatUserName(t *testing.T) {
	tests := []struct {
		name     string
		user     domain.User
		expected string
	}{
		{
			name: "Полное имя с username",
			user: domain.User{
				FirstName: "Иван",
				LastName:  "Петров",
				Username:  "ivan_petrov",
			},
			expected: "Иван Петров (@ivan_petrov)",
		},
		{
			name: "Только имя с username",
			user: domain.User{
				FirstName: "Анна",
				Username:  "anna",
			},
			expected: "Анна (@anna)",
		},
		{
			name: "Только имя без username",
			user: domain.User{
				FirstName: "Сергей",
			},
			expected: "Сергей",
		},
		{
			name: "Полное имя без username",
			user: domain.User{
				FirstName: "Мария",
				LastName:  "Сидорова",
			},
			expected: "Мария Сидорова",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.DisplayName()
			if result != tt.expected {
				t.Errorf("Ожидалось '%s', получено '%s'", tt.expected, result)
			}
		})
	}
}
