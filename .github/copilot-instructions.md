# Инструкции для GitHub Copilot

## Обзор проекта
Это Telegram бот "Пидор дня", построенный на **Go + telebot.v3** с применением принципов **чистой архитектуры**. Бот случайно выбирает "человека дня" в Telegram группах с сохранением в SQLite и комплексной системой шаблонизации.

## Ключевые архитектурные паттерны

### Слои чистой архитектуры
- **Domain** (`internal/domain/`): Чистые бизнес-сущности (`User`, `PersonOfTheDay`, `UserStats`)
- **Repository** (`internal/repository/`): Доступ к данным через **Squirrel query builder** + SQLite
- **Handlers** (`internal/handlers/`): Обработка сообщений и команд Telegram
- **Templates** (`internal/templates/`): Генерация сообщений на основе **fasttemplate**
- **Bot** (`internal/bot/`): Основной слой оркестрации

### Паттерн внедрения зависимостей
Следуйте паттерну внедрения через конструктор в `main.go`:
```go
// Сначала создаём репозитории
userRepo := repository.NewUserRepository(db)
personOfTheDayRepo := repository.NewPersonOfTheDayRepository(db)

// Затем сервисы
messageService, _ := templates.NewMessageService()

// И наконец бот со всеми зависимостями
bot := bot.NewBot(api, userRepo, personOfTheDayRepo, messageService)
```

### Паттерн интерфейсов репозиториев
Весь доступ к данным происходит через интерфейсы в `internal/repository/interfaces.go`. Конкретные реализации размещайте в отдельных файлах (`user_repository.go`, `person_of_the_day_repository.go`).

## Критические паттерны разработки

### Обработка ошибок Telegram API
Бот включает **специализированную обработку ошибок Telegram**. При добавлении отправки новых сообщений используйте паттерны обработки ошибок из `internal/handlers/utils.go`:
- `TOPIC_CLOSED` → резервная отправка в основной чат
- `chat not found`, `bot was blocked by the user` → плавная деградация
- `message is too long` → логика обрезания

### Использование системы шаблонов
**Никогда не хардкодьте пользовательские сообщения**. Весь текст должен проходить через систему шаблонов:
```go
// ✅ Правильно
message := messageService.PersonSelected(user)
// ❌ Неправильно
message := "Пидор дня: " + user.DisplayName()
```

Шаблоны используют синтаксис `{{переменная}}` с `fasttemplate`. Смотрите `internal/templates/messages.go` для всех доступных шаблонов.

### Схема базы данных и Squirrel
Используйте **Squirrel query builder** для всех SQL операций. Таблицы:
- `users`: Данные пользователей Telegram с ограничением `UNIQUE(id, chat_id)`
- `person_of_the_day`: Ежедневные выборы с отслеживанием дат

Пример паттерна запроса:
```go
query := db.psql.Select("*").From("users").Where(squirrel.Eq{"chat_id": chatID})
```

## Рабочий процесс разработки

### Настройка окружения
Необходимые переменные окружения в `.env`:
- `BOT_TOKEN`: Токен от @BotFather
- `DB_PATH`: Путь к файлу SQLite (по умолчанию: `bot.db`)
- `DEBUG`: Булево значение для режима отладки

### Команды сборки и запуска
```bash
make deps     # Установить зависимости
make build    # Собрать бинарный файл
make run      # Запустить через go run
make test     # Запустить тесты
```

### Разработка с Docker
Используйте `docker-compose.yml` для контейнеризованной разработки. Dockerfile использует многоэтапную сборку с включённым CGO для SQLite.

### Интеграция с VS Code
Проект включает `.vscode/launch.json` для отладки с переменными окружения, загружаемыми из `.env`.

## Соглашения кода

### Регистрация обработчиков
Регистрируйте обработчики в `internal/handlers/message.go` используя паттерны telebot.v3:
```go
bot.Handle("/command", handler.HandleCommand)
bot.Handle(telebot.OnText, handler.HandleMessage)
```

### Генерация случайных чисел
Используйте общий RNG из `internal/bot/bot.go` через `GetRNG()` для консистентного seeding.

### Тестирование примеров команд
Используйте `cmd/example/main.go` для тестирования вывода шаблонов без запуска полного бота.

## Точки интеграции

### Миграция на Telebot.v3
Проект мигрирован с `go-telegram-bot-api` на `gopkg.in/telebot.v3`. Используйте контекстные обработчики и паттерны middleware.

### SQLite + CGO
Требуется CGO_ENABLED=1 для `github.com/mattn/go-sqlite3`. Docker сборка корректно обрабатывает это.

### Логика работы только в группах
Бот принудительно работает только в группах. Обработчики приватных сообщений должны возвращать соответствующие шаблоны ошибок.

## Руководство по модификации файлов

### Добавление новых команд
1. Добавить метод обработчика в `internal/handlers/command.go`
2. Зарегистрировать в `internal/handlers/message.go`
3. Добавить соответствующие шаблоны в `internal/templates/messages.go`
4. Обновить шаблон текста справки

### Изменения базы данных
1. Изменить схему в `internal/repository/database.go`
2. Обновить доменные модели в `internal/domain/`
3. Добавить методы репозитория с Squirrel запросами
4. Обновить интерфейсы в `internal/repository/interfaces.go`

### Изменения шаблонов
Все шаблоны находятся в `internal/templates/` с русским текстом. Используйте синтаксис `{{переменная}}` fasttemplate.