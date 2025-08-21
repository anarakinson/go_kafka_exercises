package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/repository/database"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

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
	// Формируем DSN строку
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Println(dsn)
	// создаем соединение с базой данных
	db, err := database.NewDatabase(dsn)
	if err != nil {
		logger.Log.Error("Failed to start database", zap.Error(err))
		return
	}

	// делаем миграции
	log.Println("Starting migrations")
	if err := db.Migrate(); err != nil {
		logger.Log.Error("Migration failed, rolling back...", zap.Error(err))
		// в случае ошибки - откатываем миграции
		if rollbackErr := db.RollbackLastMigration(); rollbackErr != nil {
			logger.Log.Error("Rollback failed", zap.Error(err))
			return
		}

	}

	log.Println("Migrations successful")

}
