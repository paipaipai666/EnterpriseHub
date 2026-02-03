package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id        string         `json:"id" gorm:"primarykey"`
	OrderID   string         `json:"order_id"`
	UserID    string         `json:"user_id"`
	Amount    float64        `json:"amount"`
	Status    PaymentStatus  `json:"status"`
	Channel   PaymentChannel `json:"channel"`
	CreatedAt time.Time      `json:"create_at" gorm:"autoCreateTime"`
	PaidAt    *time.Time     `json:"paid_at,omitempty"`
}

func NewPayment(orderId, userId, channel string, amount float64) *Payment {
	return &Payment{
		Id:        uuid.New().String(),
		OrderID:   orderId,
		UserID:    userId,
		Amount:    amount,
		Status:    PaymentInit,
		Channel:   TakeChannel(channel),
		CreatedAt: time.Now(),
		PaidAt:    nil,
	}
}

func (p *Payment) Process() error {
	if p.Status != PaymentInit {
		return errors.New("invalid payment state")
	}
	p.Status = PaymentProcessing
	return nil
}

func (p *Payment) Succeed() error {
	if p.Status != PaymentProcessing {
		return errors.New("payment not processing")
	}
	p.Status = PaymentSuccess
	now := time.Now()
	p.PaidAt = &now
	return nil
}

func (p *Payment) Cancel() error {
	if p.Status != PaymentInit {
		return errors.New("invalid payment state")
	}
	p.Status = PaymentFailed
	return nil
}
