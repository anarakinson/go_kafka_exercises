# E-Commerce analytics 

Сервис, реализующий паттерн продюсер-консюмер с передачей через Kafka.

## Установка и запуск

1. Установите зависимости в каждом из сервисов (`producer`, `consumer`):
```bash
go mod tidy
```

2. Топик в Kafka создается автоматически при запуске приложения `producer`

3. Запустите сервис:

```bash
cd producer && go build -o /producer ./cmd/main.go && ./producer
cd consumer && go build -o /consumer ./cmd/main.go && ./consumer
```

4. Работа системы
Продюсер производит данные и отправляет их в кафку. Консюмер забирает накопившиеся сообщения из топика кафки

5. Как альтернатива просто используйте `docker-compose` (сначала запустите `infrastructure/docker-compose`, затем `producer/docker-compose` и `consumer/docker-compose`) 

