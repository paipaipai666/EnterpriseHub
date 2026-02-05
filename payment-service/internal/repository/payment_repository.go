package repository

import (
	"context"

	"github.com/paipaipai666/EnterpriseHub/payment-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/cache"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	Find(id string) (*domain.Payment, error)
	Save(payment *domain.Payment) error
	FindByOrderId(orderId string) ([]domain.Payment, error)
}

type paymentRepositoryImpl struct {
	cache *cache.PaymentCache
}

func NewPaymentRepository() PaymentRepository {
	return &paymentRepositoryImpl{
		cache: cache.NewPaymentCache(),
	}
}

func (pri *paymentRepositoryImpl) Create(payment *domain.Payment) error {
	err := initializers.DB.Create(&payment).Error
	if err != nil {
		return err
	}

	// 缓存
	err = pri.cache.SetPayment(context.Background(), payment)

	return err
}

func (pri *paymentRepositoryImpl) Find(id string) (*domain.Payment, error) {
	// 查缓存
	payment, err := pri.cache.GetPaymentById(context.Background(), id)
	if err != nil {
		return nil, err
	}
	if payment != nil {
		return payment, nil
	}

	payment = &domain.Payment{Id: id}
	err = initializers.DB.Find(&payment).Error
	if err != nil {
		return nil, err
	}

	return payment, nil

}

func (pri *paymentRepositoryImpl) Save(payment *domain.Payment) error {
	err := initializers.DB.Save(&payment).Error
	if err != nil {
		return err
	}

	// 缓存
	err = pri.cache.SetPayment(context.Background(), payment)

	return err
}

func (pri *paymentRepositoryImpl) FindByOrderId(orderId string) ([]domain.Payment, error) {
	payments := []domain.Payment{}
	err := initializers.DB.Where("order_id = ?", orderId).Find(&payments).Error
	if err != nil {
		return nil, err
	}

	return payments, nil
}
