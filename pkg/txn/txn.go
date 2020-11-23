// Package txn provides an interface to run data store transaction.
package txn

import (
	"context"
)

type txnContextKeyType string

// CtxKey is a context key for txn object inside a context.
const CtxKey txnContextKeyType = "ContextKey"

// Transactor is an interface for running a database transaction.
//
// It expects to create a new ctx (txnCtx) with the transaction object value inside it from the old ctx.
// Or in other words, it expects to implement something like this:
//
// txnCtx := ctx.WithValue(ctx, txn.CtxKey, <txn object>).
//
type Transactor interface {
	RunInTransaction(ctx context.Context, f func(txnCtx context.Context) error) error
}
