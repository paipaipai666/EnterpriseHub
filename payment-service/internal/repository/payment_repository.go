package repository

import (
	"github.com/paipaipai666/EnterpriseHub/payment-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	Find(id string) (*domain.Payment, error)
	Save(payment *domain.Payment) error
	FindByOrderId(orderId string) ([]domain.Payment, error)
}

type paymentRepositoryImpl struct{}

func NewPaymentRepository() PaymentRepository {
	return &paymentRepositoryImpl{}
}

func (pri *paymentRepositoryImpl) Create(payment *domain.Payment) error {
	err := initializers.DB.Create(&payment).Error

	return err
}

func (pri *paymentRepositoryImpl) Find(id string) (*domain.Payment, error) {
	payment := &domain.Payment{Id: id}
	err := initializers.DB.Find(&payment).Error
	if err != nil {
		return nil, err
	}

	return payment, nil

}

func (pri *paymentRepositoryImpl) Save(payment *domain.Payment) error {
	err := initializers.DB.Save(&payment).Error

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
