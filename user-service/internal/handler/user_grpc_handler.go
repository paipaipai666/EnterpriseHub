package handler

import (
	"context"

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
	user := ugh.repo.FindByUsername(req.Username)

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
