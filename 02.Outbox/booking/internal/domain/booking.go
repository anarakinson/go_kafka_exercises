package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/anarakinson/go_kafka_exercises/02.outbox/booking/internal/app/dto"
	"github.com/google/uuid"
)

type BookingEvent struct {
	BookingID string `json:"booking_id"`
	UserID    string `json:"user_id"`
	EventID   string `json:"event_id"`
	Timestamp int64  `json:"timestamp"`
}

type OutboxEvent struct {
	EventID   string    `json:"event_id"`   // event_id
	EventType string    `json:"event_type"` // event_type (BookingCreated)
	Payload   []byte    `json:"payload"`    //payload (JSON)
	CratedAt  time.Time `json:"created_at"` // created_at
}

func NewBookingEvent(userID string, eventID string, timestamp int64) *BookingEvent {
	return &BookingEvent{
		BookingID: uuid.New().String(),
		UserID:    userID,
		EventID:   eventID,
		Timestamp: timestamp,
	}
}

func NewOutboxEvent(eventID, eventType string, payload []byte) *OutboxEvent {
	return &OutboxEvent{
		EventID:   eventID,
		EventType: eventType,
		Payload:   payload,
		CratedAt:  time.Now(),
	}
}

// вспомогательная функция для формирования событий из дто
func MakeEvents(bReq dto.BookingRequest) (*BookingEvent, *OutboxEvent, error) {
	eventID := uuid.New()
	bookingEvent := NewBookingEvent(bReq.UserID.String(), eventID.String(), bReq.Timestamp)
	payload, err := json.Marshal(bookingEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	outboxEvent := NewOutboxEvent(eventID.String(), bReq.EventType, payload)

	return bookingEvent, outboxEvent, nil

}
