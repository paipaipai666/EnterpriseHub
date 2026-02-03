package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/domain"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/repository"
)

type PaymentGrpcHandler struct {
	pb.UnimplementedPaymentServiceServer
	repo repository.PaymentRepository
}

func NewPaymentGrpcHandler(repo repository.PaymentRepository) *PaymentGrpcHandler {
	return &PaymentGrpcHandler{
		repo: repo,
	}
}

func (pgh *PaymentGrpcHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	oldPayments, err := pgh.repo.FindByOrderId(req.OrderId)
	if err != nil {
		return nil, err
	}
	for _, oldPayment := range oldPayments {
		if oldPayment.Status == domain.PaymentSuccess {
			return nil, errors.New("请勿重复支付")
		}
	}
	payment := domain.NewPayment(req.OrderId, req.UserId, req.Method.String(), req.Amount)
	fmt.Printf("method in payment_service: %v", req.Method.String())
	err = pgh.repo.Create(payment)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentResponse{
		PaymentId: payment.Id,
		Status:    pb.PaymentStatus(payment.Status.ToStatus()),
		CreatedAt: payment.CreatedAt.Unix(),
	}, nil
}

func (pgh *PaymentGrpcHandler) Pay(ctx context.Context, req *pb.PayRequest) (*pb.PayResponse, error) {
	payment, err := pgh.repo.Find(req.PaymentId)
	res := &pb.PayResponse{
		Status:        pb.PaymentStatus(payment.Status.ToStatus()),
		TransactionId: "",
		PaidAt:        time.Time{}.Unix(),
	}

	if err != nil {
		return res, err
	}

	err = payment.Process()
	if err != nil {
		return res, err
	}

	err = payment.Succeed()
	if err != nil {
		res.Status = pb.PaymentStatus(payment.Status.ToStatus())
		return res, err
	}

	fmt.Println(payment.Status)
	pgh.repo.Save(payment)
	res.Status = pb.PaymentStatus(payment.Status.ToStatus())
	res.PaidAt = payment.PaidAt.Unix()
	return res, nil
}

func (pgh *PaymentGrpcHandler) QueryPayment(ctx context.Context, req *pb.QueryPaymentRequest) (*pb.QueryPaymentResponse, error) {
	payment, err := pgh.repo.Find(req.PaymentId)
	if err != nil {
		return nil, err
	}

	return &pb.QueryPaymentResponse{
		PaymentId: payment.Id,
		OrderId:   payment.OrderID,
		UserId:    payment.UserID,
		Amount:    payment.Amount,
		Method:    pb.PaymentMethod(payment.Channel.ToMethod()),
		Status:    pb.PaymentStatus(payment.Status.ToStatus()),
		CreatedAt: payment.CreatedAt.Unix(),
		PaidAt:    payment.PaidAt.Unix(),
	}, nil
}
