package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/anarakinson/go_kafka_exercises/01.producer_consumer/producer/internal/app/producer"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	// создаем продюсера kafka

	// Конфигурация продюсера
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Ждем подтверждения от всех реплик
	config.Producer.Retry.Max = 5                    // 5 попыток повторной отправки
	config.Producer.Return.Successes = true          // Возвращать информацию об успешной отправке

	prod := producer.NewProducer([]string{os.Getenv("KAFKA_BROKERS")}, config)
	// создаем топик, если еще не существует
	prod.CreateTopic("user-events")

	// создаем контекст
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// генерация пользовательских событий
	go func() {
		err := prod.Produce(ctx, "user-events", 3*time.Second)
		if err != nil && !errors.Is(err, producer.ErrCanceled) {
			logger.Log.Error(
				"Producer failed",
				zap.Error(err),
			)
			errChan <- err
			return
		} else {
			logger.Log.Info("Producer stopped")
			return
		}
	}()

	//--------------------------------------------//
	// создаем сервер для сборки метрик
	metric_port := os.Getenv("METRICS_PORT")
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricServer := &http.Server{
		Addr:    metric_port,
		Handler: metricsMux,
	}
	go func() {
		logger.Log.Info("Starting metrics server", zap.String("port", metric_port))
		if err := metricServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("Metrics server error", zap.Error(err))
		}
	}()

	//--------------------------------------------//
	// gracefull shutdown
	select {
	case err := <-errChan:
		logger.Log.Error("Producer service error", zap.Error(err))
	case <-shutdown:
		logger.Log.Info("Gracefully shutting down")
		// останавливаем метрик сервер
		if err := metricServer.Shutdown(ctx); err != nil {
			logger.Log.Error("Error shutting down metrics server", zap.Error(err))
		}
		// отменяем контекст
		cancel()
	}
}
