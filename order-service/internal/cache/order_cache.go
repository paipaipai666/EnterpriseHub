package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paipaipai666/EnterpriseHub/common/cache"
	"github.com/paipaipai666/EnterpriseHub/order-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/order-service/internal/domain"
)

const (
	ServerPrefix  = "enterprise_hub:order_service"
	OrderIdPrefix = "order:id:"
	UserIdPrefix  = "order:user_id:"
	OrderCacheTTL = 30 * 24 * time.Hour
	OrderLocalTTL = 24 * time.Hour
	OrderLocalMax = 1000
)

type OrderCache struct {
	redis *cache.RedisCache
	local *cache.LocalCache
}

func NewOrderCache() *OrderCache {
	return &OrderCache{
		redis: cache.NewRedisCache(ServerPrefix, initializers.RDB),
		local: cache.NewLocalCache(OrderLocalMax, OrderLocalTTL),
	}
}

func (oc *OrderCache) GetOrderById(ctx context.Context, key string) (*domain.Order, error) {
	// 本地查询
	val, err := oc.local.Get(ctx, OrderIdPrefix+key)

	if err == nil && val != nil {
		var order domain.Order

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}

		err = json.Unmarshal(data, &order)
		if err != nil {
			return nil, err
		}
		return &order, nil
	}

	// Redis缓存
	val, err = oc.redis.Get(ctx, OrderIdPrefix+key)
	if err == nil && val != "" {
		var order domain.Order

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}
		err = json.Unmarshal(data, &order)
		if err != nil {
			return nil, err
		}
		// 回填本地缓存
		oc.local.Set(ctx, OrderIdPrefix+key, val, OrderLocalTTL)
		return &order, nil
	}

	return nil, nil
}

func (oc *OrderCache) SetOrder(ctx context.Context, order *domain.Order) error {
	val, err := json.Marshal(order)
	if err != nil {
		return err
	}

	oc.local.Set(ctx, OrderIdPrefix+order.Id, val, OrderLocalTTL)

	return oc.redis.Set(ctx, OrderIdPrefix+order.Id, string(val), OrderCacheTTL)

}

func (oc *OrderCache) InvalidateOrder(ctx context.Context, order *domain.Order) error {
	oc.local.Delete(ctx, OrderIdPrefix+order.Id)

	oc.redis.Delete(ctx, OrderIdPrefix+order.Id)

	return nil
}
