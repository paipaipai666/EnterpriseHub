package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/paipaipai666/EnterpriseHub/auth-service/internal/client"
)

var hmacSampleSecret = []byte(os.Getenv("SECRET_KEY"))

type AuthService interface {
	LoginWithJWT(ctx context.Context, username, password string) (string, error)
}

type authServiceImpl struct {
	userClient client.UserClient
}

func NewAuthService(userClient client.UserClient) AuthService {
	return &authServiceImpl{
		userClient: userClient,
	}
}

func (asi *authServiceImpl) LoginWithJWT(ctx context.Context, username, password string) (string, error) {
	user, err := asi.userClient.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if password != user.Password {
		return "", errors.New("invalid password")
	}

	token, err := asi.GenerateJWT(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (asi *authServiceImpl) GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}
	return tokenString, nil
}
