package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

func NewOrder(userId string, amount float64) *Order {
	return &Order{
		Id:       uuid.New().String(),
		UserId:   userId,
		Amount:   amount,
		Status:   OrderStatusCreated,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
}

func (order *Order) Pay() error {
	if order.Status != OrderStatusCreated {
		return errors.New("order cannot be paid in current status: " + string(order.Status))
	}

	order.Status = OrderStatusPaid
	order.UpdateAt = time.Now()
	return nil
}

func (order *Order) Complete() error {
	if order.Status != OrderStatusPaid {
		return errors.New("order cannot be completed in current status: " + string(order.Status))
	}

	order.Status = OrderStatusCompleted
	order.UpdateAt = time.Now()
	return nil
}

func (order *Order) Cancel() error {
	if order.Status != OrderStatusCreated {
		return errors.New("order cannot be cancelled in current status")
	}
	order.Status = OrderStatusCancelled
	order.UpdateAt = time.Now()
	return nil
}
