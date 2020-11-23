// Package txnsql implements Transactor interface to run sql transaction.
package txnsql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/payfazz/test-api/pkg/txn"
)

// TxFromContext gets the *sqlx.Tx value from the context.
func TxFromContext(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txn.CtxKey).(*sqlx.Tx)
	return tx, ok
}

// Transactor is an implementation of dbutil.Transactor in sql.
type Transactor struct {
	db *sqlx.DB
}

// NewTransactor creates a new sql transactor.
func NewTransactor(db *sqlx.DB) *Transactor {
	return &Transactor{
		db: db,
	}
}

// RunInTransaction runs the transaction inside the context.
func (t *Transactor) RunInTransaction(ctx context.Context, f func(context.Context) error) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to create the transaction: %v", err)
	}

	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
			panic(err)
		}
	}()

	ctx = newContextTx(ctx, tx)
	if err := f(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit the transaction: %v", err)
	}

	return nil
}

func newContextTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	ctx = context.WithValue(ctx, txn.CtxKey, tx)
	return ctx
}
