package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/auth-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/api"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/client"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("auth_service")
	initializers.ConnectToRedis()
}

func main() {
	serverAddress := "0.0.0.0:8001"

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		initializers.Log.Fatal("cannot gRPC dial server: ", zap.Error(err))
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
