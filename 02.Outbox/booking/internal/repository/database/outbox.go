package database

import (
	"context"
	"fmt"
	"time"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/domain"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

// GetPendingEvents возвращает необработанные события из outbox
func (db *DataBase) GetPendingEvents(ctx context.Context) ([]domain.OutboxEvent, error) {

	// выгружаем события из бд
	rows, err := db.connection.QueryContext(
		ctx,
		`
		SELECT event_id, event_type, payload, created_at
		FROM Outbox
		WHERE processed_at IS NULL
		ORDER BY created_at
		`,
	)
	if err != nil {
		logger.Log.Error("Error executing query", zap.String("table", "Outbox"), zap.Error(err))
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	var events []domain.OutboxEvent
	for rows.Next() {
		var event domain.OutboxEvent
		err := rows.Scan(&event.EventID, &event.EventType, &event.Payload, &event.CratedAt)
		if err != nil {
			logger.Log.Error("Error reading query rows", zap.String("table", "Outbox"), zap.Error(err))
			return nil, fmt.Errorf("error reading query rows: %w", err)
		}
		events = append(events, event)
	}

	// логируем
	logger.Log.Info("Successfully load outbox events")
	return events, nil

}

// MarkEventAsProcessed помечает событие как обработанное
func (db *DataBase) MarkEventAsProcessed(ctx context.Context, eventID string) error {
	// обновляем бд
	_, err := db.connection.ExecContext(
		ctx,
		`
		UPDATE Outbox
		SET processed_at = $1
		WHERE event_id = $2
		`,
		time.Now(), eventID,
	)
	if err != nil {
		logger.Log.Error("Error updating table", zap.String("table", "Outbox"), zap.Error(err))
		return fmt.Errorf("error updating table Outbox: %w", err)
	}

	// логируем
	logger.Log.Info("Successfully update outbox processing time")
	return nil
}
