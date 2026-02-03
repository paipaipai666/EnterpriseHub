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
	"google.golang.org/grpc"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

var (
	userRepo        repository.UserRepository = repository.NewUserRepository()
	userService     service.UserService       = service.NewUserService(userRepo)
	userController  api.UserController        = api.NewUserController(userService)
	userGrpcHandler handler.UserGrpcHandler   = *handler.NewUserGrpcHandler(userRepo)
)

func main() {
	go startGrpcServer()

	startHttpServer()
}

func startGrpcServer() {
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &userGrpcHandler)

	address := "0.0.0.0:9001"
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

func startHttpServer() {
	router := gin.Default()

	apiRoutes := router.Group("/api/v1/users")
	{
		apiRoutes.POST("/register", userController.SignUp)

		apiRoutes.GET("/get/:id", userController.FindUserById)

		apiRoutes.GET("/get_all", userController.FindAllUser)
	}

	router.Run(":8000")
}
