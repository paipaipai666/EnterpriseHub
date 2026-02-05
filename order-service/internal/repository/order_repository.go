package repository

import (
	"context"

	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/cache"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
)

type OrderRepository interface {
	Save(order *domain.Order) error
	FindById(orderId string) (*domain.Order, error)
	FindByUserId(userId string) ([]domain.Order, error)
	Updata(order *domain.Order) error
}

type orderRepositoryImpl struct {
	cache *cache.OrderCache
}

func NewOrderRepository() OrderRepository {
	return &orderRepositoryImpl{
		cache: cache.NewOrderCache(),
	}
}

func (ori *orderRepositoryImpl) Save(order *domain.Order) error {
	err := initializers.DB.Create(order).Error
	if err != nil {
		return err
	}

	// 缓存
	err = ori.cache.SetOrder(context.Background(), order)

	return err
}

func (ori *orderRepositoryImpl) FindById(orderId string) (*domain.Order, error) {
	// 查缓存
	order, err := ori.cache.GetOrderById(context.Background(), orderId)
	if err != nil {
		return nil, err
	}
	if order != nil {
		return order, nil
	}

	order = &domain.Order{Id: orderId}
	err = initializers.DB.Find(&order).Error
	if err != nil {
		return nil, err
	}

	// 写缓存
	err = ori.cache.SetOrder(context.Background(), order)
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
	if err != nil {
		return err
	}

	// 缓存
	err = ori.cache.SetOrder(context.Background(), order)

	return err
}
