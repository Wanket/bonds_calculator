package util

import "fmt"

type baseTransactionError struct {
	Err   error
	TxErr error
}

func (frb *baseTransactionError) Unwrap() error {
	return frb.Err
}

type RollbackError struct {
	baseTransactionError
}

func (frb RollbackError) Error() string {
	return fmt.Sprintf("failed to rollback transaction: %v, run transaction error: %v", frb.TxErr, frb.Err)
}

type CommitError struct {
	baseTransactionError
}

func (frb CommitError) Error() string {
	return fmt.Sprintf("failed to commit transaction: %v, run transaction error: %v", frb.TxErr, frb.Err)
}

type NeedRollbackError struct {
	Cause error
}

func (nre NeedRollbackError) Unwrap() error {
	return nre.Cause
}

func (nre NeedRollbackError) Error() string {
	return fmt.Sprintf("need to rollback transaction: %v", nre.Cause)
}
