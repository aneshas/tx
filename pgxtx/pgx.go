package pgxtx

import (
	"context"
	"github.com/aneshas/tx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBFromPool(pool *pgxpool.Pool) tx.DB {
	// We can extend these to be able to receive isolation options
	// which would then be passed to tx.Begin

	return &Pool{pool}
}

type Pool struct {
	*pgxpool.Pool
}

func (p *Pool) Begin(ctx context.Context) (tx.Transaction, error) {
	txx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &PoolTx{txx}, nil
}

type PoolTx struct {
	pgx.Tx
}

func (p *PoolTx) Commit(ctx context.Context) error {
	return p.Tx.Commit(ctx)
}

func (p *PoolTx) Rollback(ctx context.Context) error {
	return p.Tx.Rollback(ctx)
}

func ConnFrom[T pgx.Tx](ctx context.Context, or T) pgx.Tx {
	if t, ok := tx.Conn[pgx.Tx](ctx); ok {
		return t.(pgx.Tx)
	}

	return or
}
