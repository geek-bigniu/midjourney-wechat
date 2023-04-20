package db

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type RedisUtil struct {
	client *redis.Client
}

func NewRedisUtil(addr string, password string, db int) *RedisUtil {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	})
	log.Printf("初始化redis成功")
	return &RedisUtil{
		client: client,
	}
}

func (ru *RedisUtil) SetNX(key string, value interface{}, expirationTime time.Duration) (bool, error) {
	result, err := ru.client.SetNX(key, value, expirationTime).Result()

	if err != nil {
		return false, err
	}

	return result, nil
}

func (ru *RedisUtil) Close() error {
	return ru.client.Close()
}
