package pgxtxv5

import (
	"context"
	"github.com/aneshas/tx/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ tx.DB = &Pool{}

// NewDBFromPool instantiates new tx.DB *pgxpool.Pool wrapper
func NewDBFromPool(pool *pgxpool.Pool, opts ...PgxTxOption) tx.DB {
	p := Pool{
		Pool: pool,
	}

	for _, opt := range opts {
		opt(&p)
	}

	return &p
}

// Pool implements tx.DB
type Pool struct {
	*pgxpool.Pool

	txOpts pgx.TxOptions
}

// Begin begins pgx transaction
func (p *Pool) Begin(ctx context.Context) (tx.Transaction, error) {
	return p.Pool.BeginTx(ctx, p.txOpts)
}

// From returns underlying pgx.Tx from the context.
// If you need to obtain a different interface back see tx.From
func From(ctx context.Context) (pgx.Tx, bool) {
	return tx.From[pgx.Tx](ctx)
}
