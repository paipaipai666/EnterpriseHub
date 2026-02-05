package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/cache"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
)

type UserRepository interface {
	CreateUser(username, password, email string) (string, error)
	FindByParam(username, password string) *model.User
	FindByUsername(username string) (*model.User, error)
	FindById(id string) *model.User
	FindAll() ([]model.User, error)
	UpdateUser(id, username, password, email string) error
}

type userRepositoryImpl struct {
	cache *cache.UserCache
}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{
		cache: cache.NewUserCache(),
	}
}

func (uri *userRepositoryImpl) CreateUser(username, password, email string) (string, error) {
	newUser := &model.User{
		Id:       uuid.New().String(),
		Username: username,
		Password: password,
		Email:    email,
	}

	err := initializers.DB.Create(&newUser).Error
	if err != nil {
		return "", err
	}
	return newUser.Id, nil
}

func (uri *userRepositoryImpl) FindByParam(username, password string) *model.User {
	user := &model.User{}
	initializers.DB.Where(&model.User{Username: username, Password: password}).First(&user)

	return user
}

func (uri *userRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	// 查缓存
	user, err := uri.cache.GetUserByUsername(context.Background(), username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		fmt.Println("----------find in cache!----------")
		return user, nil
	}

	user = &model.User{}
	err = initializers.DB.Where(&model.User{Username: username}).First(&user).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("----------find in database!----------")
	if user != nil {
		// 写缓存
		err = uri.cache.SetUser(context.Background(), user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (uri *userRepositoryImpl) FindById(id string) *model.User {
	// 查缓存
	user, err := uri.cache.GetUserById(context.Background(), id)
	if err != nil {
		return nil
	}
	if user != nil {
		return user
	}

	// 查数据库
	user = &model.User{Id: id}
	initializers.DB.First(&user)
	if user != nil {
		// 写缓存
		err = uri.cache.SetUser(context.Background(), user)
		if err != nil {
			return nil
		}
	}

	return user
}

func (uri *userRepositoryImpl) UpdateUser(id, username, password, email string) error {
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
		return initializers.DB.Limit(0).Error // 返回一个不会执行任何操作的查询
	}

	return initializers.DB.Model(&user).Where("id = ?", id).Updates(updates).Error
}

func (uri *userRepositoryImpl) FindAll() ([]model.User, error) {
	var users []model.User
	err := initializers.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
