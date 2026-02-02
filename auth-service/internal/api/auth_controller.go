package api

import (
	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/model/request"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/service"
)

type AuthController interface {
	LoginWithJWT(c *gin.Context)
	Logout(c *gin.Context)
}

type authControllerImpl struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	return &authControllerImpl{
		service: service,
	}
}

func (aci *authControllerImpl) LoginWithJWT(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "failed",
			"data":    "invalid params:" + err.Error(),
		})
	}

	token, err := aci.service.LoginWithJWT(c, loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "failed",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data":    token,
	})
}

func (aci *authControllerImpl) Logout(c *gin.Context) {
	err := aci.service.Logout(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "failed",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data":    "你已安全退出",
	})
}
