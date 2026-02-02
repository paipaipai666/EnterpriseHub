package initializers

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectToRedis() {
	rda := os.Getenv("REDIS_ADDR")

	if rda == "" {
		log.Fatal("REDIS_ADDR 环境变量未设置")
	}
	opt, err := redis.ParseURL(rda)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	RDB = redis.NewClient(opt)

	_, err = RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis数据库连接成功")
}
