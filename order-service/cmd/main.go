package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/api"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/client"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/repository"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("order_service")
	initializers.ConnectToRedis()
	initializers.ConnectToDatabase()
	initializers.ConnectToRabbitMQ()
}

func main() {
	userServerAddress := "0.0.0.0:8001"
	paymentServerAddress := "0.0.0.0:11001"

	userConn, err := grpc.NewClient(userServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		initializers.Log.Fatal("cannot user dial server: ", zap.Error(err))
	}
	paymentConn, err := grpc.NewClient(paymentServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		initializers.Log.Fatal("cannot payment dial server: ", zap.Error(err))
	}

	userServiceClient := pb.NewUserServiceClient(userConn)
	paymentServiceClient := pb.NewPaymentServiceClient(paymentConn)

	orderRepository := repository.NewOrderRepository()
	userClient := client.NewUserClient(userServiceClient)
	paymentClient := client.NewPaymentClient(paymentServiceClient)
	orderService := service.NewOrderService(orderRepository, *paymentClient, *userClient)
	orderController := api.NewOrderController(orderService)

	router := gin.Default()

	apiGroup := router.Group("/api/v1/order")
	{
		apiGroup.POST("/create", orderController.CreateOrder)

		apiGroup.GET("/get/:id", orderController.GetOrderById)

		apiGroup.GET("/list/:user_id", orderController.GetOrderList)

		apiGroup.PUT("/cancel/:id", orderController.CancelOrder)

		apiGroup.POST("/pay/:id", orderController.PayForOrder)
	}

	router.Run(":10000")
}
