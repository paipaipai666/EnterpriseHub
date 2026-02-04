package event

import "time"

type OrderCreatedEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderPaidEvent struct {
	OrderID string    `json:"order_id"`
	UserID  string    `json:"user_id"`
	Amount  float64   `json:"amount"`
	PaidAt  time.Time `json:"paid_at"`
}
