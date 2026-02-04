package main

import (
	"github.com/paipaipai666/EnterpriseHub/notify-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/notify-service/internal/consumer"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToRabbitMQ()
}

func main() {
	var forever chan struct{}

	go consumer.StartConsumer()

	<-forever
}
