package database

import (
	"context"
	"fmt"
	"time"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/app/dto"
	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/domain"
	"github.com/anarakinson/go_stonks/stonks_shared/pkg/logger"
	"go.uber.org/zap"
)

// SaveEvent сохраняет бронирование и событие в одной транзакции
func (db *DataBase) SaveEvent(ctx context.Context, bookingRequest dto.BookingRequest) error {
	// начало транзакции
	tx, err := db.connection.Begin()
	if err != nil {
		logger.Log.Error("Error making events", zap.Error(err))
		return fmt.Errorf("transaction failed: %w", err)
	}
	// откатывваем в случае неудачи
	defer tx.Rollback()

	// создаем события
	bookingEvent, outboxEvent, err := domain.MakeEvents(bookingRequest)
	if err != nil {
		logger.Log.Error("Error making events", zap.Error(err))
		return fmt.Errorf("error making events: %w", err)
	}

	// создае6м outbox событие
	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO Outbox (event_id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4)
		`,
		outboxEvent.EventID, outboxEvent.EventType, outboxEvent.Payload, outboxEvent.CratedAt,
	)
	if err != nil {
		logger.Log.Error("Error inserting into table", zap.String("table", "Outbox"), zap.Error(err))
		return fmt.Errorf("error inserting into Outbox: %w", err)
	}

	// сохраняем бронирование
	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO Booking (booking_id, user_id, event_id, created_at)
		VALUES ($1, $2, $3, $4)
		`,
		bookingEvent.BookingID, bookingEvent.UserID, bookingEvent.EventID, time.Unix(bookingEvent.Timestamp, 0),
	)
	if err != nil {
		logger.Log.Error("Error inserting into table", zap.String("table", "Booking"), zap.Error(err))
		return fmt.Errorf("error inserting into Booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		logger.Log.Error("Error committing", zap.Error(err))
		return fmt.Errorf("error committing: %w", err)
	}

	logger.Log.Info("Database successfully updated")
	return nil

}
