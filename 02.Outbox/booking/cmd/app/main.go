package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/app/delivery"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/infrastructure"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/repository/database"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/server"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/pkg/metrics"

	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

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
	fmt.Println(dsn)
	db, err := database.NewDatabase(dsn)
	if err != nil {
		logger.Log.Error("Error connecting to database", zap.Error(err))
		return
	}
	defer db.Close()

	//--------------------------------------------//
	// создаем общий контекст для отмены всего
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//--------------------------------------------//
	// создаем нового продюсера
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	producer, err := infrastructure.NewProducer([]string{os.Getenv("KAFKA_BROKERS")}, config)
	if err != nil {
		logger.Log.Error("Error creating Kafka producer", zap.Error(err))
		return
	}

	// создаем топик в кафке
	topicName := "booking-events"
	producer.CreateTopic(topicName)

	//--------------------------------------------//
	// запускаем воркера для отправки сообщений в кафку
	worker := infrastructure.NewWorker(db, producer)
	go worker.Work(ctx, topicName, 30*time.Second)

	//--------------------------------------------//
	// создаем сервер, слушающий сообщения
	serv := server.NewServer()

	// создаем хэндлер
	handler := delivery.NewHandler(db)

	// запускаем сервер
	port := os.Getenv("PORT")
	go func() {
		err := serv.Run(port, handler)
		if err != nil && err != http.ErrServerClosed {
			errChan <- err
			logger.Log.Error("Error serving", zap.String("port", port), zap.Error(err))
			return
		}
	}()

	//--------------------------------------------//
	// создаем сервер для сборки метрик
	metric_port := os.Getenv("METRICS_PORT")
	metricServer := metrics.NewMetricServer()
	go func() {
		logger.Log.Info("Start metrics server:", zap.String("port", port))
		if err := metricServer.Run(metric_port); err != nil {
			logger.Log.Error("Metrics server error:", zap.Error(err))
		}
	}()

	//--------------------------------------------//
	// gracefull shutdown
	select {
	case err := <-errChan:
		logger.Log.Error("Main service error", zap.Error(err))
	case <-shutdown:
		logger.Log.Info("Gracefully shutting down")
		// останавливаем метрик сервер
		if err := metricServer.Shutdown(ctx); err != nil {
			logger.Log.Error("Error shutting down metrics server", zap.Error(err))
		}
		// останавливаем основной сервер
		if err := serv.Shutdown(ctx); err != nil {
			logger.Log.Error("Error shutting down main server", zap.Error(err))
		}
		// отменяем контекст
		cancel()
	}
}
