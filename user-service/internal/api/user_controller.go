package api

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model/dto"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/service"
)

type UserController interface {
	SignUp(ctx *gin.Context)
	FindUserById(ctx *gin.Context)
	FindUserByUsername(ctx *gin.Context)
	FindAllUser(ctx *gin.Context)
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

	userId, err := uci.service.SignUp(singUpData)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    userId,
	})
}

func (uci *userControllerImpl) FindUserById(ctx *gin.Context) {
	id := ctx.Param("id")

	user := uci.service.GetUserById(id)

	if user == nil {
		ctx.JSON(500, gin.H{
			"message": "failed",
			"data":    id,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    user,
	})
}

func (uci *userControllerImpl) FindUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := uci.service.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	if user == nil {
		ctx.JSON(404, gin.H{
			"message": "failed",
			"data":    "user not found",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    user,
	})
}

func (uci *userControllerImpl) FindAllUser(ctx *gin.Context) {
	users, err := uci.service.GetAllUser()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success",
		"data":    users,
	})
}
