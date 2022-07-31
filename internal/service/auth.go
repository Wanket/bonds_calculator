package service

import (
	"bonds_calculator/internal/data/repository"
	datautil "bonds_calculator/internal/data/util"
	"bonds_calculator/internal/util"
	"context"
	"errors"
	"fmt"
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword

type IAuthService interface {
	CreateUser(context context.Context, username string, password string) (Token, error)
	LoginUser(context context.Context, username string, password string) (Token, error)

	ChangePassword(ctx context.Context, username string, oldPassword string, newPassword string) error

	RefreshToken(context context.Context, token string) (Token, error)

	GetSharedSecret() []byte
}

type Token struct {
	Bytes string

	ExpirationTime time.Time
}

type AuthService struct {
	userRepository  datautil.IRepositoryBuilder[repository.IUserRepository]
	tokenRepository repository.ITokenRepository

	sharedSecret []byte

	clock clock.Clock

	refreshTokenTTL time.Duration
}

func NewAuthService(
	tokenRepository repository.ITokenRepository,
	userRepository datautil.IRepositoryBuilder[repository.IUserRepository],
	clock clock.Clock,
	config util.IGlobalConfig,
) *AuthService {
	sharedSecret, err := getOrGenerateSharedSecret(tokenRepository)
	if err != nil {
		log.Fatalf("AuthService: %v", err)
	}

	return &AuthService{
		tokenRepository: tokenRepository,
		userRepository:  userRepository,

		sharedSecret: sharedSecret,

		clock: clock,

		refreshTokenTTL: config.GetRefreshTokenTTL(),
	}
}

func (auth *AuthService) CreateUser(ctx context.Context, username string, password string) (Token, error) {
	transaction, err := auth.userRepository.WithTransaction(ctx)
	if err != nil {
		return Token{}, fmt.Errorf("failed to create transaction: %w", err)
	}

	return datautil.OnTransaction(transaction, func(repo repository.IUserRepository) (Token, error) { //nolint:wrapcheck
		exist, err := repo.IsUserExists(ctx, username)
		if err != nil {
			return Token{}, fmt.Errorf("failed to check if user exists: %w", err)
		}

		if exist {
			return Token{}, repository.ErrUserAlreadyExists
		}

		hash, err := bcrypt.GenerateFromPassword(util.StringToBytes(password), bcrypt.DefaultCost)
		if err != nil {
			return Token{}, fmt.Errorf("failed to generate password hash: %w", err)
		}

		if err = repo.CreateUserIfNotExists(ctx, username, hash); err != nil {
			if errors.Is(err, repository.ErrUserAlreadyExists) {
				return Token{}, repository.ErrUserAlreadyExists
			}

			return Token{}, fmt.Errorf("failed to create user: %w", err)
		}

		token, err := auth.generateAndStoreToken(ctx)
		if err != nil {
			return Token{}, datautil.NeedRollbackError{ //nolint:exhaustruct
				Cause: fmt.Errorf("failed to generate token: %w", err),
			}
		}

		return token, nil
	})
}

func (auth *AuthService) LoginUser(ctx context.Context, username string, password string) (Token, error) {
	repo := auth.userRepository.WithoutTransaction()

	hash, err := repo.GetUserCredentials(ctx, username)
	if err != nil {
		return Token{}, fmt.Errorf("failed to get user credentials: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(hash, util.StringToBytes(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return Token{}, ErrMismatchedHashAndPassword
		}

		return Token{}, fmt.Errorf("failed to compare password hash: %w", err)
	}

	return auth.generateAndStoreToken(ctx)
}

func (auth *AuthService) ChangePassword(
	ctx context.Context,
	username string,
	oldPassword string,
	newPassword string,
) error {
	transaction, err := auth.userRepository.WithTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	_, err = datautil.OnTransaction(
		transaction,
		func(repo repository.IUserRepository) (struct{}, error) {
			hash, err := repo.GetUserCredentials(ctx, username)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to get user credentials: %w", err)
			}

			if err = bcrypt.CompareHashAndPassword(hash, util.StringToBytes(oldPassword)); err != nil {
				if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
					return struct{}{}, ErrMismatchedHashAndPassword
				}

				return struct{}{}, fmt.Errorf("failed to compare password hash: %w", err)
			}

			hash, err = bcrypt.GenerateFromPassword(util.StringToBytes(newPassword), bcrypt.DefaultCost)
			if err != nil {
				return struct{}{}, fmt.Errorf("failed to generate password hash: %w", err)
			}

			if err = repo.UpdateUserCredentials(ctx, username, hash); err != nil {
				return struct{}{}, fmt.Errorf("failed to update user credentials: %w", err)
			}

			return struct{}{}, nil
		})

	return err //nolint:wrapcheck
}

func (auth *AuthService) RefreshToken(ctx context.Context, token string) (Token, error) {
	if err := auth.tokenRepository.DeleteToken(ctx, token); err != nil {
		if errors.Is(err, repository.ErrTokenNotFoundOrExpired) {
			return Token{}, repository.ErrTokenNotFoundOrExpired
		}

		return Token{}, fmt.Errorf("failed to delete token: %w", err)
	}

	return auth.generateAndStoreToken(ctx)
}

func (auth *AuthService) GetSharedSecret() []byte {
	return auth.sharedSecret
}

func (auth *AuthService) generateAndStoreToken(ctx context.Context) (Token, error) {
	tokenBytes, err := util.GenerateSecretString(util.TokenSize)
	if err != nil {
		return Token{}, fmt.Errorf("failed to generate token: %w", err)
	}

	token := Token{
		Bytes:          tokenBytes,
		ExpirationTime: auth.clock.Now().Add(auth.refreshTokenTTL),
	}

	if err := auth.tokenRepository.StoreToken(ctx, token.Bytes, auth.refreshTokenTTL); err != nil {
		return Token{}, fmt.Errorf("failed to store token: %w", err)
	}

	return token, nil
}

func getOrGenerateSharedSecret(tokenRepository repository.ITokenRepository) ([]byte, error) {
	sharedSecret, err := util.GenerateSecret(util.SharedSecretSize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate shared secret: %w", err)
	}

	sharedSecret, err = tokenRepository.SetNXAndGetSharedSecret(sharedSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get shared secret: %w", err)
	}

	return sharedSecret, nil
}
