package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paipaipai666/EnterpriseHub/common/cache"
	"github.com/paipaipai666/EnterpriseHub/user-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/user-service/internal/model"
)

const (
	ServerPrefix   = "enterprise_hub:user_service"
	UserIdPrefix   = "user:id:"
	UserNamePrefix = "user:username:"
	UserCacheTTL   = 30 * 24 * time.Hour
	UserLocalTTL   = 5 * time.Minute
	UserLocalMax   = 1000
)

type UserCache struct {
	redis *cache.RedisCache
	local *cache.LocalCache
}

func NewUserCache() *UserCache {
	return &UserCache{
		redis: cache.NewRedisCache(ServerPrefix, initializers.RDB),
		local: cache.NewLocalCache(UserLocalMax, UserLocalTTL),
	}
}

func (uc *UserCache) GetUserByUsername(ctx context.Context, key string) (*model.User, error) {
	// 本地查询
	val, err := uc.local.Get(ctx, UserNamePrefix+key)
	if err == nil && val != nil {
		var user model.User

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}
		err = json.Unmarshal(data, &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	// Redis缓存
	val, err = uc.redis.Get(ctx, UserNamePrefix+key)
	if err == nil && val != "" {
		var user model.User

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}
		err = json.Unmarshal(data, &user)
		if err != nil {
			return nil, err
		}
		// 回填本地缓存
		uc.local.Set(ctx, UserNamePrefix+key, val, UserLocalTTL)
		return &user, nil
	}

	return nil, nil
}

func (uc *UserCache) GetUserById(ctx context.Context, key string) (*model.User, error) {
	// 本地查询
	val, err := uc.local.Get(ctx, UserIdPrefix+key)
	if err == nil && val != nil {
		var user model.User

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}
		err = json.Unmarshal(data, &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}

	// Redis缓存
	val, err = uc.redis.Get(ctx, UserIdPrefix+key)
	if err == nil && val != "" {
		var user model.User

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}
		err = json.Unmarshal(data, &user)
		if err != nil {
			return nil, err
		}
		// 回填本地缓存
		uc.local.Set(ctx, UserIdPrefix+key, val, UserLocalTTL)
		return &user, nil
	}

	return nil, nil
}

func (uc *UserCache) SetUser(ctx context.Context, user *model.User) error {
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}

	uc.local.Set(ctx, UserIdPrefix+user.Id, val, UserLocalTTL)
	uc.local.Set(ctx, UserNamePrefix+user.Username, val, UserLocalTTL)

	err = uc.redis.Set(ctx, UserIdPrefix+user.Id, string(val), UserCacheTTL)
	if err != nil {
		return err
	}
	err = uc.redis.Set(ctx, UserNamePrefix+user.Username, string(val), UserCacheTTL)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserCache) InvalidateUser(ctx context.Context, user *model.User) error {
	uc.local.Delete(ctx, UserIdPrefix+user.Id)

	uc.redis.Delete(ctx, UserIdPrefix+user.Id)

	return nil
}
