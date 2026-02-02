package service

import (
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
	dto "github.com/paipaipai666/EnterpriseHub/user-service/internal/model/dto"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(param *dto.UserDTO) *gorm.DB
	GetUserByUsername(username string) *model.User
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

func (usi *userServiceImpl) SignUp(param *dto.UserDTO) *gorm.DB {
	return usi.repo.CreateUser(param.Username, param.Password, param.Email)
}

func (usi *userServiceImpl) GetUserByUsername(username string) *model.User {
	user := usi.repo.FindByUsername(username)

	return user
}

func (usi *userServiceImpl) GetUserById(id string) *model.User {
	user := usi.repo.FindById(id)

	return user
}

func (usi *userServiceImpl) GetAllUser() ([]model.User, error) {
	users, result := usi.repo.FindAll()
	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}
