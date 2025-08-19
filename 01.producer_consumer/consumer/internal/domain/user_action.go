package domain

type UserEventType string

var (
	EventView      UserEventType = "view"
	EventAddToCart UserEventType = "add_to_cart"
	EventPurchase  UserEventType = "purchase"
)

type UserAction struct {
	UserID    string        `json:"user_id"`
	ProductID string        `json:"product_id"`
	EventType UserEventType `json:"event_type"`
	Timestamp int64         `json:"timestamp"`
}
