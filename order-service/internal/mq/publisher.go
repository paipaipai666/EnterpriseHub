package mq

import (
	"encoding/json"

	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishOrderCreated(order *domain.Order) error {
	event := event.OrderCreatedEvent{
		OrderID:   order.Id,
		UserID:    order.UserId,
		Amount:    order.Amount,
		CreatedAt: order.CreateAt,
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return initializers.Channel.Publish("order_events", "order.created", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func PublishOrderPaid(order *domain.Order) error {
	event := event.OrderPaidEvent{
		OrderID: order.Id,
		UserID:  order.UserId,
		Amount:  order.Amount,
		PaidAt:  order.UpdateAt,
	}
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return initializers.Channel.Publish("order_events", "order.paid", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}
