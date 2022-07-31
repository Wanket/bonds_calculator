package connection

import (
	"bonds_calculator/internal/util"
	"fmt"
	"github.com/go-redis/redis/v9"
)

type RedisConnection struct {
	*redis.Client
}

func NewRedisConnection(config util.IGlobalConfig) *RedisConnection {
	redisConfig := config.GetRedisConfig()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	return &RedisConnection{
		Client: redisClient,
	}
}
