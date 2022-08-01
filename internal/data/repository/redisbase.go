package repository

import (
	"bonds_calculator/internal/data/connection"
	"fmt"
)

type BaseRedisRepository struct {
	conn *connection.RedisConnection
}

func NewBaseRedisRepository(conn *connection.RedisConnection) BaseRedisRepository {
	return BaseRedisRepository{
		conn: conn,
	}
}

func (repository *BaseRedisRepository) Close() error {
	if err := repository.conn.Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}

	return nil
}
