# Используем официальный образ Go
FROM golang:1.24 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем исходный код
COPY . .

# Скачиваем зависимости
RUN go mod download

# собираем бинарник, делающий миграции
RUN CGO_ENABLED=0 GOOS=linux go build -o /migrate ./cmd/migrations/main.go
# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /entrypoint ./cmd/app/main.go

# Финальный образ
FROM alpine:latest

WORKDIR /app

COPY .env .env
COPY ./migrations ./migrations/

COPY --from=builder /migrate ./
COPY --from=builder /entrypoint ./

# Запускаем приложение
CMD ["sh", "-c", "./migrate && tail -f /app/logs/app.log & ./entrypoint > /app/logs/app.log 2>&1"]
