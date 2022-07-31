package repository

import (
	"bonds_calculator/internal/data/connection"
	"context"
	"fmt"
	"time"
)

const (
	keyStorePrefix  = "user.token.store."
	sharedSecretKey = "user.token.shared_secret"
)

var ErrTokenNotFoundOrExpired = fmt.Errorf("token not found or expired")

type ITokenRepository interface {
	StoreToken(ctx context.Context, key string, expirationTime time.Duration) error
	DeleteToken(ctx context.Context, key string) error

	SetNXAndGetSharedSecret(sharedSecret []byte) ([]byte, error)
}

type TokenRepository struct {
	BaseRedisRepository
}

func NewTokenRepository(conn *connection.RedisConnection) *TokenRepository {
	return &TokenRepository{
		BaseRedisRepository: NewBaseRedisRepository(conn),
	}
}

func (repository *TokenRepository) StoreToken(ctx context.Context, key string, expirationTime time.Duration) error {
	err := repository.conn.Set(ctx, keyStorePrefix+key, 0, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	return nil
}

func (repository *TokenRepository) DeleteToken(ctx context.Context, key string) error {
	result, err := repository.conn.Del(ctx, keyStorePrefix+key).Result()
	if result == 0 {
		return ErrTokenNotFoundOrExpired
	}

	if err != nil {
		return fmt.Errorf("failed to delete token: %w", err)
	}

	return nil
}

func (repository *TokenRepository) SetNXAndGetSharedSecret(sharedSecret []byte) ([]byte, error) {
	err := repository.conn.SetNX(context.Background(), sharedSecretKey, sharedSecret, 0).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set shared secret: %w", err)
	}

	result, err := repository.conn.Get(context.Background(), sharedSecretKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %w", err)
	}

	return result, nil
}
