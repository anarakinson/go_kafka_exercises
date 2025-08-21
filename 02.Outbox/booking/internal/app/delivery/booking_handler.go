package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/app/dto"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	SaveEvent(context.Context, dto.BookingRequest) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) BookingHandler(w http.ResponseWriter, r *http.Request) {

	// ироверяем, поддерживается ли метод
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method is not supported"}`, http.StatusMethodNotAllowed)
		return
	}

	// декодируем тело запроса
	var bReq dto.BookingRequest
	err := json.NewDecoder(r.Body).Decode(&bReq)
	if err != nil {
		http.Error(w, `{"error":"Wrong json format"}`, http.StatusBadRequest)
		return
	}

	// валидация полей
	if bReq.UserID == uuid.Nil || bReq.EventType == "" || bReq.Timestamp == 0 {
		http.Error(w, `{"error":"Required fields: [user_id, event_type, timestamp]"}`, http.StatusBadRequest)
		return
	}

	// логируем
	logger.Log.Info(
		"message received",
		zap.String("user_id", bReq.UserID.String()),
		zap.String("event_type", bReq.EventType),
		zap.Int64("timestamp", bReq.Timestamp),
	)

	// создаем бронирование в базе данных
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = h.repo.SaveEvent(ctx, bReq)
	if err != nil {
		logger.Log.Error("Error updating database", zap.Error(err))
		http.Error(w, `{"error":"Error saving data"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Crated",
	})

}
