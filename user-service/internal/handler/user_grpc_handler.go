package handler

import (
	"context"
	"fmt"

	"github.com/paipaipai666/EnterpriseHub/user-service/internal/pb"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/repository"
)

type UserGrpcHandler struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
}

func NewUserGrpcHandler(repo repository.UserRepository) *UserGrpcHandler {
	return &UserGrpcHandler{repo: repo}
}

func (ugh *UserGrpcHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.UserResponse, error) {
	user, err := ugh.repo.FindByUsername(req.Username)
	if err != nil {
		return nil, err
	}

	fmt.Printf("查询到的用户: ID=%s, Username=%s\n", user.Id, user.Username)

	return &pb.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}, nil
}

func (ugh *UserGrpcHandler) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.UserResponse, error) {
	user := ugh.repo.FindById(req.Id)

	return &pb.UserResponse{
		Id:       user.Id,
		Username: user.Username,
		Password: user.Password,
	}, nil
}
