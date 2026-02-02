package service

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/client"
)

var hmacSampleSecret = []byte(os.Getenv("SECRET_KEY"))

type AuthService interface {
	LoginWithJWT(ctx *gin.Context, username, password string) (string, error)
	Logout(ctx *gin.Context) error
}

type authServiceImpl struct {
	userClient   client.UserClient
	tokenChecker UserTokenService
}

func NewAuthService(userClient client.UserClient, tokenChecker UserTokenService) AuthService {
	return &authServiceImpl{
		userClient:   userClient,
		tokenChecker: tokenChecker,
	}
}

func (asi *authServiceImpl) LoginWithJWT(ctx *gin.Context, username, password string) (string, error) {
	user, err := asi.userClient.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if password != user.Password {
		return "", errors.New("invalid password")
	}

	err = asi.tokenChecker.TokenExists(username, ctx)
	if err != nil {
		return "", err
	}

	token, err := asi.tokenChecker.GenerateJWT(username)
	if err != nil {
		return "", err
	}

	asi.tokenChecker.SaveToken(username, token, ctx)
	return token, nil
}

func (asi *authServiceImpl) Logout(ctx *gin.Context) error {
	err := asi.tokenChecker.DeleteToken(ctx)

	return err
}
