package tx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// BeginFunc executes the provided function within a database transaction.
//
// If the function returns an error, the transaction is rolled back, otherwise it is committed.
func BeginFunc(
	ctx context.Context,
	db *sql.DB,
	fn func(tx *sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil && !errors.Is(err, sql.ErrConnDone) {
			err = rollbackErr
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
