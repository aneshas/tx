package sqltx

import (
	"context"
	"database/sql"
	"github.com/aneshas/tx/v2"
)

var (
	_ tx.DB          = &DB{}
	_ tx.Transaction = &Tx{}
)

// NewDB instantiates new tx.DB *sql.DB wrapper
func NewDB(db *sql.DB) tx.DB {
	return &DB{db}
}

// DB implements tx.DB
type DB struct {
	*sql.DB
}

// Begin begins sql transaction
func (p *DB) Begin(_ context.Context) (tx.Transaction, error) {
	txx, err := p.DB.Begin()
	if err != nil {
		return nil, err
	}

	return &Tx{txx}, nil
}

// Tx wraps *sql.TX in order top implement tx.Transaction
type Tx struct {
	*sql.Tx
}

// Commit commits the transaction
func (p *Tx) Commit(_ context.Context) error {
	return p.Tx.Commit()
}

// Rollback rolls back the transaction
func (p *Tx) Rollback(_ context.Context) error {
	return p.Tx.Rollback()
}

// From returns underlying *sql.Tx (wrapped in *Tx)
func From(ctx context.Context) (*Tx, bool) {
	return tx.From[*Tx](ctx)
}
