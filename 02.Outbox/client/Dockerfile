# Используем официальный образ Go
FROM golang:1.24 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем исходный код
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint ./cmd/app/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY .env .env
COPY --from=builder /entrypoint ./

# Запускаем приложение
# CMD ["./entrypoint"]
# CMD ["sh", "-c", "./entrypoint > /app/logs/app.log 2>&1"]
CMD ["sh", "-c", "tail -f /app/logs/app.log & ./entrypoint > /app/logs/app.log 2>&1"]

