package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/api"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/repository"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/service"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

var (
	orderRepository repository.OrderRepository = repository.NewOrderRepository()
	orderService    service.OrderService       = service.NewOrderService(orderRepository)
	orderController api.OrderController        = api.NewOrderController(orderService)
)

func main() {
	router := gin.Default()

	apiGroup := router.Group("/api/v1/order")
	{
		apiGroup.POST("/create", orderController.CreateOrder)

		apiGroup.GET("/get/:id", orderController.GetOrderById)

		apiGroup.GET("/list/:user_id", orderController.GetOrderList)

		apiGroup.PUT("/cancel/:id", orderController.CancelOrder)
	}

	router.Run(":10000")
}
