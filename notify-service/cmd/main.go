package main

import (
	"github.com/paipaipai666/EnterpriseHub/notify-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/notify-service/internal/consumer"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("notify_service")
	initializers.ConnectToRabbitMQ()
}

func main() {
	forever := make(chan struct{})

	go consumer.StartConsumer()

	<-forever
}
