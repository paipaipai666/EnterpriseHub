package initializers

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RDB *redis.Client

func ConnectToRedis() {
	rda := os.Getenv("REDIS_ADDR")

	if rda == "" {
		Log.Fatal("REDIS_ADDR 环境变量未设置")
	}
	opt, err := redis.ParseURL(rda)
	if err != nil {
		Log.Fatal("连接数据库失败: %v", zap.Error(err))
	}

	RDB = redis.NewClient(opt)

	_, err = RDB.Ping(context.Background()).Result()
	if err != nil {
		Log.Fatal("Failed to connect to Redis: %v", zap.Error(err))
	}
	Log.Info("Redis数据库连接成功")
}
