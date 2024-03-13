package boiltx

import (
	"context"
	"database/sql"
	"github.com/aneshas/tx"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func NewDB(db *sql.DB) tx.DB {
	return &DB{db}
}

type DB struct {
	*sql.DB
}

func (p *DB) Begin(_ context.Context) (tx.Transaction, error) {
	txx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &PoolTx{txx}, nil
}

type PoolTx struct {
	*sql.Tx
}

func (p *PoolTx) Commit(_ context.Context) error {
	return p.Tx.Commit()
}

func (p *PoolTx) Rollback(_ context.Context) error {
	return p.Tx.Rollback()
}

func foo() {
	t := tx.New(NewDB(nil))

	_ = t.Do(context.TODO(), func(ctx context.Context) error {
		return nil
	})
}

func ConnFrom[T boil.ContextExecutor](ctx context.Context, or T) boil.ContextExecutor {
	if t, ok := tx.Conn[*sql.Tx](ctx); ok {
		return t.(*sql.Tx)
	}

	return or
}
