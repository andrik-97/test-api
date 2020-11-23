// Package txnmock provides mock implementation of Transactor to use in testing.
package txnmock

import (
	"context"

	"github.com/payfazz/test-api/pkg/txn"
)

// Tx is a mock transaction object.
type Tx struct{}

// TxnFromContext gets the *Tx value from the context.
func TxnFromContext(ctx context.Context) (*Tx, bool) {
	tx, ok := ctx.Value(txn.CtxKey).(*Tx)
	return tx, ok
}

// Transactor is a mock implementation of txn.Transactor.
type Transactor struct {
	RunInTransactionInvoked bool
}

// NewTransactor creates a new mock transactor.
func NewTransactor() *Transactor {
	return &Transactor{}
}

// RunInTransaction runs the transaction inside the context.
func (t *Transactor) RunInTransaction(ctx context.Context, f func(context.Context) error) error {
	t.RunInTransactionInvoked = true

	ctx = newContextTxn(ctx, &Tx{})
	if err := f(ctx); err != nil {
		return err
	}

	return nil
}

func newContextTxn(ctx context.Context, tx *Tx) context.Context {
	ctx = context.WithValue(ctx, txn.CtxKey, tx)
	return ctx
}
