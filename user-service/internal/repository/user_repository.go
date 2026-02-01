package repository

import (
	"github.com/google/uuid"
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(username, password, email string) *gorm.DB
	FindByParam(username, password string) *model.User
	FindById(id string) *model.User
	UpdateUser(id, username, password, email string) *gorm.DB
}

type userRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (uri *userRepositoryImpl) CreateUser(username, password, email string) *gorm.DB {
	newUser := &model.User{
		Id:       uuid.New().String(),
		Username: username,
		Password: password,
		Email:    email,
	}

	return initializers.DB.Create(&newUser)
}

func (uri *userRepositoryImpl) FindByParam(username, password string) *model.User {
	user := &model.User{}
	initializers.DB.Where(&model.User{Username: username, Password: password}).First(&user)

	return user
}

func (uri *userRepositoryImpl) FindById(id string) *model.User {
	user := &model.User{Id: id}
	initializers.DB.First(&user)

	return user
}

func (uri *userRepositoryImpl) UpdateUser(id, username, password, email string) *gorm.DB {
	user := &model.User{}
	var updates = make(map[string]interface{})

	if username != "" {
		updates["username"] = username
	}
	if password != "" {
		updates["password"] = password
	}
	if email != "" {
		updates["email"] = email
	}

	if len(updates) == 0 {
		return initializers.DB.Limit(0) // 返回一个不会执行任何操作的查询
	}

	return initializers.DB.Model(&user).Where("id = ?", id).Updates(updates)
}
