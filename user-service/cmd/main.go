package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/api"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/handler"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/repository"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/service"
	"github.com/paipaipai666/EnterpriseHub/user-service/middleware"
	"google.golang.org/grpc"
)

func init() {
	initializers.LoadEnv()
	initializers.InitLogger("user_service")
	initializers.ConnectToRedis()
	initializers.ConnectToDatabase()
}

func main() {
	// 在所有初始化完成后才创建服务实例
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userController := api.NewUserController(userService)
	userGrpcHandler := *handler.NewUserGrpcHandler(userRepo)

	go startGrpcServer(userGrpcHandler)

	startHttpServer(userController)
}

func startGrpcServer(userGrpcHandler handler.UserGrpcHandler) {
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userGrpcHandler)

	address := "0.0.0.0:8001"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func startHttpServer(userController api.UserController) {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.GinLogger)

	apiRoutes := router.Group("/api/v1/users")
	{
		apiRoutes.POST("/register", userController.SignUp)

		apiRoutes.GET("/get/:id", userController.FindUserById)

		apiRoutes.GET("/get_by_username/:username", userController.FindUserByUsername)

		apiRoutes.GET("/get_all", userController.FindAllUser)
	}

	router.Run(":8000")
}
