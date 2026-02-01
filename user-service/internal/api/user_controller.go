package api

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model/dto"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/service"
)

type UserController interface {
	SignUp(ctx *gin.Context)
}

type userControllerImpl struct {
	service service.UserService
}

func NewUserController(service service.UserService) UserController {
	return &userControllerImpl{
		service: service,
	}
}

func (uci *userControllerImpl) SignUp(ctx *gin.Context) {
	singUpData := &dto.UserDTO{}

	err := ctx.ShouldBindJSON(&singUpData)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	result := uci.service.SignUp(singUpData)
	if result.Error != nil {
		ctx.JSON(500, gin.H{
			"message": "failed",
			"data":    result.Error.Error(),
		})
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    "user_id",
	})
}
