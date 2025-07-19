package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/Martinpasaribu/Golang-V1/internal/config"
)

func HasViewedBlog(blogID string, ip string) bool {
	key := fmt.Sprintf("viewed:%s:%s", blogID, ip)
	exists, err := config.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	return exists == 1
}

func SetViewedBlog(blogID string, ip string, duration time.Duration) {
	key := fmt.Sprintf("viewed:%s:%s", blogID, ip)
	config.RedisClient.Set(context.Background(), key, "1", duration)
}
