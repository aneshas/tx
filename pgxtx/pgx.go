package pgxtx

import (
	"context"
	"github.com/aneshas/tx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ tx.DB = &Pool{}

// NewDBFromPool instantiates new tx.DB *pgxpool.Pool wrapper
func NewDBFromPool(pool *pgxpool.Pool) tx.DB {
	// We can extend these to be able to receive isolation options
	// which would then be passed to tx.Begin

	return &Pool{pool}
}

// Pool represents tx wrapper for *pgxpool.Pool in order to implement tx.DB
type Pool struct {
	*pgxpool.Pool
}

// Begin begins pgx transaction
func (p *Pool) Begin(ctx context.Context) (tx.Transaction, error) {
	return p.Pool.Begin(ctx)
}

// From returns underlying T from the context which in most cases should probably be pgx.Tx
// but is left as a generic type in order to accommodate cases where people tend to abstract
// the whole connection away behind an interface (see examples)
// T should be an interface
//
// Example:
func From[T any](ctx context.Context, pool *pgxpool.Pool) (T, error) {
	return tx.From[T](ctx, NewDBFromPool(pool))
}
