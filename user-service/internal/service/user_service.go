package service

import (
	dto "github.com/paipaipai666/EnterpriseHub/user-service/internal/model/dto"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/repository"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(param *dto.UserDTO) *gorm.DB
	Login(param *dto.LoginDTO) string
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

func (usi *userServiceImpl) Login(param *dto.LoginDTO) string {
	user := usi.repo.FindByParam(param.Username, param.Password)
	if user == nil {
		return ""
	}

	return user.Id
}
