package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatalf("❌ Gagal parse Redis URL: %v", err)
	}

	RedisClient = redis.NewClient(opt)

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Gagal konek Redis: %v", err)
	}

	log.Println("✅ Redis connected (serverless)")
}
