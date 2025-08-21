# Booking Service

Сервис для обработки бронирований с отправкой событий в Kafka.

## Установка и запуск

1. Установите зависимости:
```bash
cd booking && go mod tidy
```

2. Создайте базу данных в PostgreSQL. Запустите миграции в директории `migrations`, либо с помощью `cmd/migrations/main.go`

3. Топик в Kafka создается автоматически при запуске приложения

4. Запустите сервис:

```bash
cd booking && go build -o /booking ./cmd/app/main.go && ./booking
```

5. Использование
Отправьте POST запрос для создания бронирования:

bash
```bash
curl -X POST http://localhost:8080/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "b2c3d4e5-f6g7-8901-bcde-fg2345678901", 
    "event_type": "BookingCreated",
    "timestamp": "2024-01-15T10:30:00Z"
  }'
```

6. Как альтернатива просто используйте `docker-compose` (сначала запустите `infrastructure/docker-compose`, затем `booking/docker-compose`). Для отправки сообщений можно использовать `client/docker-compose`
