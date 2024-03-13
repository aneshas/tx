package sqltx

import (
	"context"
	"database/sql"
	"github.com/aneshas/tx"
)

var (
	_ tx.DB          = &DB{}
	_ tx.Transaction = &Tx{}
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

	return &Tx{txx}, nil
}

type Tx struct {
	*sql.Tx
}

func (p *Tx) Commit(_ context.Context) error {
	return p.Tx.Commit()
}

func (p *Tx) Rollback(_ context.Context) error {
	return p.Tx.Rollback()
}

func From(ctx context.Context, db *sql.DB) (*Tx, error) {
	return tx.From[*Tx](ctx, NewDB(db))
}
