package client

import (
	"context"

	"github.com/paipaipai666/EnterpriseHub/order-service/internal/pb"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentClient(conn pb.PaymentServiceClient) *PaymentClient {
	return &PaymentClient{client: conn}
}

func (pc *PaymentClient) CreatePayment(ctx context.Context, orderId, userId string, amount float64, method int) (*pb.CreatePaymentResponse, error) {
	return pc.client.CreatePayment(ctx, &pb.CreatePaymentRequest{
		OrderId: orderId,
		UserId:  userId,
		Amount:  amount,
		Method:  pb.PaymentMethod(method),
	})
}

func (pc *PaymentClient) Pay(ctx context.Context, paymentId string) (*pb.PayResponse, error) {
	return pc.client.Pay(ctx, &pb.PayRequest{
		PaymentId: paymentId,
	})
}

func (pc *PaymentClient) QueryPayment(ctx context.Context, req *pb.QueryPaymentRequest) (*pb.QueryPaymentResponse, error) {
	return pc.client.QueryPayment(ctx, req)
}
