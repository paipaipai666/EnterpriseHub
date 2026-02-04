package initializers

import (
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

func ConnectToRabbitMQ() {
	url := os.Getenv("RABBITMQ_URL")

	var err error
	Conn, err = amqp.Dial(url)
	if err != nil {
		Log.Fatal("RabbitMQ 连接失败: ", zap.Error(err))
	}

	Channel, err = Conn.Channel()
	if err != nil {
		Log.Fatal("RabbitMQ 通道创建失败: ", zap.Error(err))
	}

	// 声明 exchange 和 queue
	Channel.ExchangeDeclare("order_events", "topic", true, false, false, false, nil)
	Channel.QueueDeclare("notifications", true, false, false, false, nil)
	Channel.QueueBind("notifications", "order.#", "order_events", false, nil)
}
