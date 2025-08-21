package dto

import "github.com/google/uuid"

type BookingRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	EventType string    `json:"event_type"`
	Timestamp int64     `json:"timestamp"`
}
