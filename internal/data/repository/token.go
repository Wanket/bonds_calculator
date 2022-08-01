package repository

import (
	"bonds_calculator/internal/data/connection"
	"bonds_calculator/internal/data/entgenerated"
	"bonds_calculator/internal/data/entgenerated/secret"
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

	SetNXAndGetSharedSecret(ctx context.Context, sharedSecret []byte) ([]byte, error)
}

type TokenRepository struct {
	BaseRedisRepository

	dbConn *connection.DBConnection
}

func NewTokenRepository(conn *connection.RedisConnection, dbConnection *connection.DBConnection) *TokenRepository {
	return &TokenRepository{
		BaseRedisRepository: NewBaseRedisRepository(conn),
		dbConn:              dbConnection,
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

func (repository *TokenRepository) SetNXAndGetSharedSecret(ctx context.Context, sharedSecret []byte) ([]byte, error) {
	err := repository.dbConn.Secret.Create().SetKey(sharedSecretKey).SetValue(sharedSecret).Exec(ctx)
	if err == nil {
		return sharedSecret, nil
	}

	if !entgenerated.IsConstraintError(err) {
		return nil, fmt.Errorf("failed to set shared secret: %w", err)
	}

	result, err := repository.dbConn.Secret.Query().Where(secret.Key(sharedSecretKey)).Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %w", err)
	}

	return result.Value, nil
}
