package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/paipaipai666/EnterpriseHub/common/cache"
	"github.com/paipaipai666/EnterpriseHub/payment-service/initializers"
	"github.com/paipaipai666/EnterpriseHub/payment-service/internal/domain"
)

const (
	ServerPrefix    = "enterprise_hub:payment_service"
	PaymentIdPrefix = "payment:id:"
	UserIdPrefix    = "payment:user_id:"
	PaymentCacheTTL = 30 * 24 * time.Hour
	PaymentLocalTTL = 24 * time.Hour
	PaymentLocalMax = 1000
)

type PaymentCache struct {
	redis *cache.RedisCache
	local *cache.LocalCache
}

func NewPaymentCache() *PaymentCache {
	return &PaymentCache{
		redis: cache.NewRedisCache(ServerPrefix, initializers.RDB),
		local: cache.NewLocalCache(PaymentLocalMax, PaymentLocalTTL),
	}
}

func (pc *PaymentCache) GetPaymentById(ctx context.Context, key string) (*domain.Payment, error) {
	// 本地查询
	val, err := pc.local.Get(ctx, PaymentIdPrefix+key)

	if err == nil && val != nil {
		var payment domain.Payment

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}

		err = json.Unmarshal(data, &payment)
		if err != nil {
			return nil, err
		}
		return &payment, nil
	}

	// Redis缓存
	val, err = pc.redis.Get(ctx, PaymentIdPrefix+key)
	if err == nil && val != "" {
		var payment domain.Payment

		var data []byte
		switch v := val.(type) {
		case string:
			data = []byte(v)
		case []byte:
			data = v
		default:
			return nil, fmt.Errorf("cached value is not of expected type (string or []byte), got %T", v)
		}
		err = json.Unmarshal(data, &payment)
		if err != nil {
			return nil, err
		}
		// 回填本地缓存
		pc.local.Set(ctx, PaymentIdPrefix+key, val, PaymentLocalTTL)
		return &payment, nil
	}

	return nil, nil
}

func (pc *PaymentCache) SetPayment(ctx context.Context, payment *domain.Payment) error {
	val, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	pc.local.Set(ctx, PaymentIdPrefix+payment.Id, val, PaymentLocalTTL)

	return pc.redis.Set(ctx, PaymentIdPrefix+payment.Id, string(val), PaymentCacheTTL)

}

func (pc *PaymentCache) InvalidatePayment(ctx context.Context, payment *domain.Payment) error {
	pc.local.Delete(ctx, PaymentIdPrefix+payment.Id)

	pc.redis.Delete(ctx, PaymentIdPrefix+payment.Id)

	return nil
}
