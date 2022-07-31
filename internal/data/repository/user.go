package repository

import (
	"bonds_calculator/internal/data/connection"
	"bonds_calculator/internal/data/entgenerated"
	User "bonds_calculator/internal/data/entgenerated/user"
	datautil "bonds_calculator/internal/data/util"
	"bonds_calculator/internal/util"
	"context"
	"fmt"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserNotFound      = fmt.Errorf("user not found")
)

type IUserRepository interface {
	IsUserExists(ctx context.Context, username string) (bool, error)
	CreateUserIfNotExists(ctx context.Context, username string, passwordHash []byte) error
	GetUserCredentials(ctx context.Context, username string) (passwordHash []byte, err error)
	UpdateUserCredentials(ctx context.Context, username string, passwordHash []byte) error
}

type UserRepository struct {
	BaseDBRepository
}

func NewUserRepository(conn *connection.DBConnection) *UserRepository {
	return &UserRepository{
		BaseDBRepository: NewBaseDBRepository(conn),
	}
}

func (repository *UserRepository) WithTransaction(ctx context.Context) (datautil.Transaction[IUserRepository], error) {
	return withTransaction[IUserRepository](ctx, util.CopyObjectByRef(repository))
}

func (repository *UserRepository) WithoutTransaction() IUserRepository { //nolint:ireturn
	return repository
}

func (repository *UserRepository) IsUserExists(ctx context.Context, username string) (bool, error) {
	user := repository.getUserTable()

	_, err := user.Query().Where(User.UserName(username)).Only(ctx)
	if err != nil {
		if entgenerated.IsNotFound(err) {
			return false, nil
		}

		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return true, nil
}

func (repository *UserRepository) CreateUserIfNotExists(
	ctx context.Context,
	username string,
	passwordHash []byte,
) error {
	user := repository.getUserTable()

	_, err := user.Create().SetUserName(username).SetPasswordHash(passwordHash).Save(ctx)
	if entgenerated.IsConstraintError(err) {
		return ErrUserAlreadyExists
	}

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (repository *UserRepository) GetUserCredentials(
	ctx context.Context,
	username string,
) ([]byte, error) {
	user := repository.getUserTable()

	userEntity, err := user.Query().Where(User.UserName(username)).Only(ctx)
	if err != nil {
		if entgenerated.IsNotFound(err) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user credentials: %w", err)
	}

	return userEntity.PasswordHash, nil
}

func (repository *UserRepository) UpdateUserCredentials(
	ctx context.Context,
	username string,
	passwordHash []byte,
) error {
	user := repository.getUserTable()

	count, err := user.Update().Where(User.UserName(username)).SetPasswordHash(passwordHash).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user credentials: %w", err)
	}

	if count == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (repository *UserRepository) getUserTable() *entgenerated.UserClient {
	if repository.tx != nil {
		return repository.tx.User
	}

	return repository.conn.User
}
