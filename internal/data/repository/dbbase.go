package repository

import (
	"bonds_calculator/internal/data/connection"
	"bonds_calculator/internal/data/entgenerated"
	"bonds_calculator/internal/data/util"
	"context"
	"database/sql/driver"
	"fmt"
)

type Repo[IT any] interface {
	WithoutTransaction() IT

	CreateAndSetTx(ctx context.Context) (driver.Tx, error)
}

type BaseDBRepository struct {
	conn *connection.DBConnection

	tx *entgenerated.Tx
}

func NewBaseDBRepository(conn *connection.DBConnection) BaseDBRepository {
	return BaseDBRepository{
		conn: conn,

		tx: nil,
	}
}

func (repository *BaseDBRepository) CreateAndSetTx(ctx context.Context) (driver.Tx, error) {
	transaction, err := repository.conn.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	repository.tx = transaction

	return transaction, nil
}

func (repository *BaseDBRepository) Close() error {
	if err := repository.conn.Close(); err != nil {
		return fmt.Errorf("failed to close db connection: %w", err)
	}

	return nil
}

func withTransaction[IT any](ctx context.Context, repository Repo[IT]) (util.Transaction[IT], error) {
	transaction, err := repository.CreateAndSetTx(ctx)
	if err != nil {
		return util.Transaction[IT]{}, err //nolint:wrapcheck
	}

	return util.Transaction[IT]{
		// We set transaction, so WithoutTransaction just return interface using transaction
		Repository: repository.WithoutTransaction(),
		Tx:         transaction,
	}, nil
}
