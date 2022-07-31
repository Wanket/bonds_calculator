package util

import (
	"context"
	"database/sql/driver"
	"errors"
	log "github.com/sirupsen/logrus"
)

type IRepositoryBuilder[Repository any] interface {
	WithTransaction(ctx context.Context) (Transaction[Repository], error)

	WithoutTransaction() Repository
}

type Transaction[Repository any] struct {
	Repository Repository

	Tx driver.Tx
}

func OnTransaction[Repository any, R any, F func(Repository) (result R, err error)]( //nolint:nolintlint,ireturn
	transaction Transaction[Repository],
	onTransactionF F,
) (R, error) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}

		log.WithFields(
			log.Fields{
				"panic err":    err,
				"rollback err": transaction.Tx.Rollback(),
			},
		).Error("got panic during transaction")

		panic(err)
	}()

	var (
		nilResult         R
		needRollbackError NeedRollbackError
	)

	result, err := onTransactionF(transaction.Repository)
	if errors.As(err, &needRollbackError) {
		if rbErr := transaction.Tx.Rollback(); rbErr != nil {
			return nilResult, &RollbackError{
				baseTransactionError: baseTransactionError{
					Err:   needRollbackError.Unwrap(),
					TxErr: rbErr,
				},
			}
		}

		return result, needRollbackError.Unwrap()
	}

	if cErr := transaction.Tx.Commit(); cErr != nil {
		return nilResult, &CommitError{
			baseTransactionError: baseTransactionError{
				Err:   err,
				TxErr: cErr,
			},
		}
	}

	return result, err
}
