package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/api"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/repository"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/service"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDatabase()
}

var (
	userRepo       repository.UserRepository = repository.NewUserRepository()
	userService    service.UserService       = service.NewUserService(userRepo)
	userController api.UserController        = api.NewUserController(userService)
)

func main() {
	router := gin.Default()

	apiRoutes := router.Group("/api/v1")
	{
		apiRoutes.POST("/register", userController.SignUp)

		apiRoutes.GET("/login", userController.Login)
	}

	router.Run(":8080")
}
