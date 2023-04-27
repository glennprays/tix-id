package tool

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func SetRedisValue(client *redis.Client, key, value string, expiration time.Duration) error {
	err := client.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisValue(client *redis.Client, key string) (string, error) {
	value, err := client.Get(key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key does not exist")
	} else if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
	return client
}
