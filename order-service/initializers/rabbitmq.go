package initializers

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection
var Channel *amqp.Channel

func ConnectToRabbitMQ() {
	url := os.Getenv("RABBITMQ_URL")

	var err error
	Conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatal("RabbitMQ 连接失败: " + err.Error())
	}

	Channel, err = Conn.Channel()
	if err != nil {
		log.Fatal("RabbitMQ 通道创建失败: " + err.Error())
	}

	// 声明 exchange 和 queue
	Channel.ExchangeDeclare("order_events", "topic", true, false, false, false, nil)
	Channel.QueueDeclare("notifications", true, false, false, false, nil)
	Channel.QueueBind("notifications", "order.#", "order_events", false, nil)
}
