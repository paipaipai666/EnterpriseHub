package consumer

import (
	"encoding/json"
	"log"

	"github.com/paipaipai666/EnterpriseHub/notify-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/notify-service/internal/event"
	"github.com/paipaipai666/EnterpriseHub/notify-service/internal/service"
)

func StartConsumer() {
	msgs, _ := initializers.Channel.Consume("notifications", "", false, false, false, false, nil)

	for msg := range msgs {
		// 使用消息的RoutingKey来判断消息类型
		routingKey := msg.RoutingKey

		switch routingKey {
		case "order.created":
			var event event.OrderCreatedEvent
			json.Unmarshal(msg.Body, &event)

			notificationService := service.NewNotificationService()
			notificationService.SendOrderCreatedNotification(event)
		case "order.paid":
			var event event.OrderPaidEvent
			json.Unmarshal(msg.Body, &event)

			notificationService := service.NewNotificationService()
			notificationService.SendOrderPaidNotification(event)
		default:
			// 处理未知消息类型
			log.Printf("未知的消息类型: %s", routingKey)
		}

		msg.Ack(false)
	}
}
