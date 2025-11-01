# Используем официальный образ Go для сборки
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты для сборки SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY *.go ./

# Собираем бота
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bot .

# Финальный образ
FROM alpine:latest

# Устанавливаем sqlite и ca-certificates
RUN apk --no-cache add ca-certificates sqlite

# Создаем пользователя для бота
RUN adduser -D -s /bin/sh bot

# Создаем директорию для приложения
WORKDIR /app

# Копируем бинарник из builder образа
COPY --from=builder /app/bot .

# Создаем директорию для базы данных
RUN mkdir -p /app/data && chown bot:bot /app/data

# Переключаемся на пользователя bot
USER bot

# Указываем переменные окружения по умолчанию
ENV DB_PATH=/app/data/bot.db
ENV DEBUG=false

# Открываем порт (не обязательно для Telegram бота, но может пригодиться)
EXPOSE 8080

# Запускаем бота
CMD ["./bot"]