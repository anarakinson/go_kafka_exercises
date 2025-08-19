package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"itk/lessons/07.kafka/02.Outbox/booking/database"

	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type BookingRequest struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	EventID   string `json:"event_id"`
	Timestamp int32  `json:"timestamp"`
}

type OutboxEvent struct {
	EventID   string    // event_id
	EventType string    // event_type (BookingCreated)
	Payload   []byte    //payload (JSON)
	CratedAt  time.Time // created_at
}

func main() {

	//--------------------------------------------//
	// Канал для graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	errChan := make(chan error, 1)

	//--------------------------------------------//
	// загружаем переменные окружения
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		return
	}

	//--------------------------------------------//
	// инициализируем логгер
	if err := logger.Init("production"); err != nil {
		slog.Error("Unable to init zap-logger", "error", err)
		return
	}
	defer logger.Sync()

	//--------------------------------------------//
	// создаем соединение с базой данных
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := database.NewDatabase(dsn)
	if err != nil {
		logger.Log.Error("Error connecting to database", zap.Error(err))
	}

	//--------------------------------------------//
	// создаем топик в кафке

	//--------------------------------------------//
	// gracefull shutdown
	select {
	case err := <-errChan:
		logger.Log.Error("Producer service error", zap.Error(err))
	case <-shutdown:
		logger.Log.Info("Gracefully shutting down")
		// останавливаем метрик сервер
		// if err := metricServer.Shutdown(ctx); err != nil {
		// 	logger.Log.Error("Error shutting down metrics server", zap.Error(err))
		// }
		// // отменяем контекст
		// cancel()
	}
}
