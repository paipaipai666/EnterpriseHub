package service

import (
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(userId string, amount float64) (string, error)
	GetOrder(id string) (*domain.Order, error)
	GetOrderList(userId string) ([]domain.Order, error)
	Cancel(id string) error
}

type orderServiceImpl struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &orderServiceImpl{
		repo: repo,
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
