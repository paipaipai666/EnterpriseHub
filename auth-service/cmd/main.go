package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/auth-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/api"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/client"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToRedis()
}

func main() {
	serverAddress := "0.0.0.0:9001"

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	userServiceClient := pb.NewUserServiceClient(conn)
	userClient := client.NewUserClient(userServiceClient)
	tokenChecker := service.NewUserTokenService()
	authService := service.NewAuthService(*userClient, tokenChecker)
	authController := api.NewAuthController(authService)

	router := gin.Default()

	apiRoutes := router.Group("/api/v1")
	{
		apiRoutes.POST("/auth/login", authController.LoginWithJWT)

		apiRoutes.DELETE("/auth/logout", authController.Logout)
	}

	router.Run(":9000")
}
