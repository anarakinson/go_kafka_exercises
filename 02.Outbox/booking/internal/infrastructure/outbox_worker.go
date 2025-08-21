package infrastructure

import (
	"context"
	"time"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/domain"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

type Repository interface {
	GetPendingEvents(context.Context) ([]domain.OutboxEvent, error)
	MarkEventAsProcessed(context.Context, string) error
}

type Worker struct {
	repo     Repository
	producer *Producer
}

func NewWorker(repo Repository, producer *Producer) *Worker {
	return &Worker{
		repo:     repo,
		producer: producer,
	}
}

func (w *Worker) Work(ctx context.Context, topicName string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop worker")
			return
		case <-ticker.C:
			w.processOutboxEvents(ctx, topicName)
		}
	}
}

func (w *Worker) processOutboxEvents(ctx context.Context, topicName string) {

	// даем три секунды на обработку
	ctx, _ = context.WithTimeout(ctx, 3*time.Second)
	events, err := w.repo.GetPendingEvents(ctx)
	if err != nil {
		logger.Log.Error("Error getting events", zap.Error(err))
		return
	}

	if len(events) == 0 {
		logger.Log.Info("Nothing to do")
		return
	}

	logger.Log.Info("Found unprocessed events", zap.Int("events_count", len(events)))

	counter := 0
	for _, event := range events {
		// отправляем в кафку с гарантией "at least once"
		err := w.producer.SendToKafka(topicName, event)
		if err != nil {
			logger.Log.Error("Error sending event to Kafka", zap.Error(err))
			continue
		}

		// даем три секунды на обработку
		ctx, _ = context.WithTimeout(ctx, 3*time.Second)
		// обновляем базу данных
		err = w.repo.MarkEventAsProcessed(ctx, event.EventID)
		if err != nil {
			logger.Log.Error("Error marking event as processed in database", zap.Error(err))
			continue
		}
		// считаем, сколько упешно обработано
		counter++
	}

	logger.Log.Info("Events successfully processed", zap.Int("processed_events", counter), zap.Int("total_events", len(events)))

}
