package main

import (
	"fmt"
	"time"

	"github.com/pavel-one/day-of-the-bot/internal/domain"
	"github.com/pavel-one/day-of-the-bot/internal/templates"
)

func main() {
	service, err := templates.NewMessageService()
	if err != nil {
		fmt.Printf("Ошибка создания сервиса: %v\n", err)
		return
	}

	// Создаем тестового пользователя
	user := domain.User{
		ID:        123456,
		FirstName: "Павел",
		LastName:  "Зарубин",
		Username:  "pavelzarubin",
	}

	fmt.Println("=== Примеры сообщений бота ===")
	fmt.Println()

	fmt.Println("1. Справка:")
	fmt.Println(service.HelpText())
	fmt.Println()

	fmt.Println("2. Пидор дня выбран:")
	fmt.Println(service.PersonSelected(user))
	fmt.Println()

	fmt.Println("3. Пидор дня уже выбран:")
	fmt.Println(service.PersonAlreadySelected(user))
	fmt.Println()

	fmt.Println("4. Информация о пидоре дня:")
	fmt.Println(service.PersonInfo(user, time.Now()))
	fmt.Println()

	fmt.Println("5. Статистика:")
	stats := []domain.UserStats{
		{User: user, Count: 10},
		{User: domain.User{ID: 2, FirstName: "Иван", LastName: "Петров"}, Count: 7},
		{User: domain.User{ID: 3, FirstName: "Мария", Username: "maria"}, Count: 5},
	}
	fmt.Println(service.StatsHeader(stats))
	fmt.Println()

	fmt.Println("6. Нет статистики:")
	fmt.Println(service.NoStatsAvailable())
	fmt.Println()

	fmt.Println("7. Ошибка:")
	fmt.Println(service.ErrorOccurred("при подключении к базе данных"))
	fmt.Println()

	fmt.Println("8. Неизвестная команда:")
	fmt.Println(service.UnknownCommand())
}
