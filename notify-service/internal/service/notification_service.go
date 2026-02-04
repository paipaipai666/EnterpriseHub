package service

import (
	"log"

	"github.com/paipaipai666/EnterpriseHub/notify-service/internal/event"
)

type NotificationService interface {
	SendOrderCreatedNotification(event event.OrderCreatedEvent) error
	SendOrderPaidNotification(event event.OrderPaidEvent) error
}

type notificationServiceImpl struct{}

func NewNotificationService() NotificationService {
	return &notificationServiceImpl{}
}

func (nsi *notificationServiceImpl) SendOrderCreatedNotification(event event.OrderCreatedEvent) error {
	log.Printf("[通知] 订单创建成功 - 订单ID: %s, 用户ID: %s, 金额: %.2f, 创建时间: %v",
		event.OrderID, event.UserID, event.Amount, event.CreatedAt)
	return nil
}

func (nsi *notificationServiceImpl) SendOrderPaidNotification(event event.OrderPaidEvent) error {
	log.Printf("[通知] 订单支付成功 - 订单ID: %s, 用户ID: %s, 金额: %.2f, 支付时间: %v",
		event.OrderID, event.UserID, event.Amount, event.PaidAt)
	return nil
}
