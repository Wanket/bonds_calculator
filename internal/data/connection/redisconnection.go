package connection

import (
	"bonds_calculator/internal/util"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	log "github.com/sirupsen/logrus"
	"time"
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

	_, err := redisClient.Ping(context.Background()).Result()
	for err != nil {
		log.WithError(err).Error("RedisConnection: Failed to connect to redis, retrying...")

		time.Sleep(time.Second)

		_, err = redisClient.Ping(context.Background()).Result()
	}

	log.Info("RedisConnection: Redis connection established")

	return &RedisConnection{
		Client: redisClient,
	}
}
