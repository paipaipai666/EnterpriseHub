package cache

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	rda := os.Getenv("REDIS_ADDR")
	opt, err := redis.ParseURL(rda)
	if err != nil {
		panic(err)
	}
	RedisClient = redis.NewClient(opt)
}

type RedisCache struct {
	client *redis.Client
	prefix string
}

func NewRedisCache(prefix string, redisClient *redis.Client) *RedisCache {
	return &RedisCache{
		client: redisClient,
		prefix: prefix,
	}
}

func (r *RedisCache) key(key string) string {
	return r.prefix + ":" + key
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.key(key)).Result()
}

func (r *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, r.key(key), value, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.key(key)).Err()
}

func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, r.key(key)).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
