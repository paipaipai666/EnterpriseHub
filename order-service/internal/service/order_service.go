package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/client"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(userId string, amount float64) (string, error)
	Pay(ctx *gin.Context, method int, orderId string) (string, error)
	GetOrder(id string) (*domain.Order, error)
	GetOrderList(userId string) ([]domain.Order, error)
	Cancel(id string) error
}

type orderServiceImpl struct {
	repo          repository.OrderRepository
	paymentClient client.PaymentClient
	userClient    client.UserClient
}

func NewOrderService(repo repository.OrderRepository, paymentClient client.PaymentClient, userClient client.UserClient) OrderService {
	return &orderServiceImpl{
		repo:          repo,
		paymentClient: paymentClient,
		userClient:    userClient,
	}
}

func (osi *orderServiceImpl) CreateOrder(userId string, amount float64) (string, error) {
	order := domain.NewOrder(userId, amount)
	err := osi.repo.Save(order)
	if err != nil {
		return "", err
	}

	return order.Id, nil
}

func (osi *orderServiceImpl) Pay(ctx *gin.Context, method int, orderId string) (string, error) {
	username, err := osi.getUsername(ctx)
	if err != nil {
		return "", fmt.Errorf("无法获取用户名: %v", err) // 返回更具体的错误信息
	}

	userRes, err := osi.userClient.GetUserByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("无法获取用户信息: %v", err)
	}

	order, err := osi.repo.FindById(orderId)
	if err != nil {
		return "", fmt.Errorf("无法找到订单: %v", err)
	}

	paymentRes, err := osi.paymentClient.CreatePayment(ctx, order.Id, userRes.Id, order.Amount, method)
	if err != nil {
		return "", fmt.Errorf("创建支付失败: %v", err)
	}

	payRes, err := osi.paymentClient.Pay(ctx, paymentRes.PaymentId)
	if err != nil {
		return "", fmt.Errorf("支付处理失败: %v", err)
	}

	switch payRes.Status {
	case 1:
		return "", errors.New("未完成支付！")
	case 3:
		return "", errors.New("支付失败！")
	case 4:
		return "", errors.New("支付超时！")
	case 0:
		return "", errors.New("未知支付错误！")
	}
	return paymentRes.PaymentId, nil
}

func (osi *orderServiceImpl) GetOrder(id string) (*domain.Order, error) {
	order, err := osi.repo.FindById(id)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (osi *orderServiceImpl) GetOrderList(userId string) ([]domain.Order, error) {
	orders, err := osi.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (osi *orderServiceImpl) Cancel(id string) error {
	order, err := osi.repo.FindById(id)
	if err != nil {
		return err
	}
	order.Cancel()

	err = osi.repo.Updata(order)
	return err
}

func (osi *orderServiceImpl) getUsername(ctx *gin.Context) (string, error) {
	var hmacSampleSecret = []byte(os.Getenv("SECRET_KEY"))

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
