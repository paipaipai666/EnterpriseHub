package repository

import (
	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
)

type OrderRepository interface {
	Save(order *domain.Order) error
	FindById(orderId string) (*domain.Order, error)
	FindByUserId(userId string) ([]domain.Order, error)
	Updata(order *domain.Order) error
}

type orderRepositoryImpl struct {
}

func NewOrderRepository() OrderRepository {
	return &orderRepositoryImpl{}
}

func (ori *orderRepositoryImpl) Save(order *domain.Order) error {
	err := initializers.DB.Create(order).Error

	return err
}

func (ori *orderRepositoryImpl) FindById(orderId string) (*domain.Order, error) {
	order := &domain.Order{Id: orderId}
	err := initializers.DB.Find(&order).Error
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (ori *orderRepositoryImpl) FindByUserId(userId string) ([]domain.Order, error) {
	orders := []domain.Order{}
	err := initializers.DB.Where("user_id = ?", userId).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, err
}

func (ori *orderRepositoryImpl) Updata(order *domain.Order) error {
	err := initializers.DB.Save(order).Error

	return err
}
