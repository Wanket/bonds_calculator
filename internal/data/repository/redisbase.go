package repository

import (
	"bonds_calculator/internal/data/connection"
)

type BaseRedisRepository struct {
	conn *connection.RedisConnection
}

func NewBaseRedisRepository(conn *connection.RedisConnection) BaseRedisRepository {
	return BaseRedisRepository{
		conn: conn,
	}
}
