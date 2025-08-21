package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/client/internal/app/client"

	"github.com/joho/godotenv"
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
	// запускаем отправку запросов на букинг
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// создаем воркера
	worker := client.NewWorker()
	if err != nil {
		slog.Error("Error starting worker", "error", err)
		return
	}

	// формируем адрес
	url := fmt.Sprintf("http://%s/bookings", os.Getenv("SERVER_ADDRESS"))
	// запускаем отправку в фоновом режиме
	slog.Info("Start working", "url", url)
	go func() {
		err := worker.Work(ctx, url)
		if err != nil {
			errChan <- err
			return
		}
	}()

	//--------------------------------------------//
	// gracefull shutdown
	select {
	case <-errChan:
		slog.Error("Error while working", "error", err)
	case <-shutdown:
		slog.Info("Shutting down")
	}

}
