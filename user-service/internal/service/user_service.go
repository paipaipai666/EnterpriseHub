package service

import (
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
	dto "github.com/paipaipai666/EnterpriseHub/user-service/internal/model/dto"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/repository"
)

type UserService interface {
	SignUp(param *dto.UserDTO) (string, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserById(id string) *model.User
	GetAllUser() ([]model.User, error)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}

func (usi *userServiceImpl) SignUp(param *dto.UserDTO) (string, error) {
	userId, err := usi.repo.CreateUser(param.Username, param.Password, param.Email)
	if err != nil {
		return "", err
	}
	return userId, nil
}

func (usi *userServiceImpl) GetUserByUsername(username string) (*model.User, error) {
	user, err := usi.repo.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (usi *userServiceImpl) GetUserById(id string) *model.User {
	user := usi.repo.FindById(id)

	return user
}

func (usi *userServiceImpl) GetAllUser() ([]model.User, error) {
	users, err := usi.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
