package client

import (
	"context"

	"github.com/paipaipai666/EnterpriseHub/order-service/internal/pb"
)

type UserClient struct {
	client pb.UserServiceClient
}

func NewUserClient(conn pb.UserServiceClient) *UserClient {
	return &UserClient{client: conn}
}
func (uc *UserClient) GetUserByUsername(ctx context.Context, username string) (*pb.GetUserResponse, error) {
	return uc.client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{
		Username: username,
	})
}
