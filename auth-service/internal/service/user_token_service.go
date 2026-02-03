package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/paipaipai666/EnterpriseHub/auth-service/initializers"
	"github.com/redis/go-redis/v9"
)

type UserTokenService interface {
	SaveToken(username, token string, ctx *gin.Context) error
	DeleteToken(ctx *gin.Context) error
	TokenExists(username string, ctx *gin.Context) error
	GenerateJWT(username string) (string, error)
}

type userTokenServiceImpl struct {
	keyHand string
}

func NewUserTokenService() UserTokenService {
	return &userTokenServiceImpl{
		keyHand: "enterprise_hub:jwt:login:",
	}
}

func (usti *userTokenServiceImpl) SaveToken(username, token string, ctx *gin.Context) error {
	return initializers.RDB.Set(ctx, usti.keyHand+username, token, 24*time.Hour).Err()
}

func (usti *userTokenServiceImpl) TokenExists(username string, ctx *gin.Context) error {
	val, err := initializers.RDB.Get(ctx, usti.keyHand+username).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	if val != "" {
		return errors.New("用户已登录")
	}
	return nil
}

func (usti *userTokenServiceImpl) GenerateJWT(username string) (string, error) {
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

func (usti *userTokenServiceImpl) DeleteToken(ctx *gin.Context) error {
	username, err := usti.getUsername(ctx)
	if err != nil {
		return err
	}

	delCount, err := initializers.RDB.Del(ctx, usti.keyHand+username).Result()
	if err != nil {
		return err
	}
	if delCount == 0 {
		return errors.New("该用户未登录")
	}
	return nil
}

func (usti *userTokenServiceImpl) getUsername(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("请求头中无 Authorization 参数")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSampleSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return "", errors.New("Failed to parse token:" + err.Error())
	}

	if !token.Valid {
		return "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("Invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		fmt.Println("Invalid or missing expiration time")
		return "", errors.New("Invalid expiration time")
	}

	if exp < float64(time.Now().Unix()) {
		return "", errors.New("Token has expired")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("属性'username'不为string类型")
	}
	return username, nil
}
