package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp int64     `json:"timestamp"`
}

func NewBookingRequest(eventType string) *BookingRequest {
	return &BookingRequest{
		UserID:    uuid.New(),
		EventType: eventType,
		Timestamp: time.Now().Unix(),
	}
}
