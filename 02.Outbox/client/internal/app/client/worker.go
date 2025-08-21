package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/client/internal/domain"
)

var ErrPingFail = errors.New("url not responding")

type BookingResponse struct {
	Status string `json:"status"`
}

type Worker struct {
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) Work(ctx context.Context, url string) error {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("context canceled")
			return nil
		case <-ticker.C:
			msg := domain.NewBookingRequest("BookingCreated")
			msgJson, err := json.Marshal(msg)
			if err != nil {
				slog.Error("Error marchalling message", "error", err)
				continue
			}
			// создаем сообщение с контекстом
			req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(msgJson))
			if err != nil {
				slog.Error("Error creating request with context", "error", err)
				continue
			}

			// отправляем сообщение на сервер
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				slog.Error("Error sending message", "error", err, "url", url)
				continue
			}
			// закрываем тело ответа
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				slog.Error("Status code not ok", "status_code", resp.StatusCode, "url", url)
				continue
			}

			// Читаем тело ответа
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				slog.Error("Cannot read response", "error", err)
				continue
			}

			// Декодируем JSON
			var bookingResp BookingResponse
			err = json.Unmarshal(bodyBytes, &bookingResp)
			if err != nil {
				slog.Error("Wrong JSON format", "error", err)
				continue
			}

			slog.Info("Response received", "status_code", resp.StatusCode, "body", bookingResp)

		}

	}
}
