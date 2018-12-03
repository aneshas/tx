package postgres

import (
	"context"
	"database/sql"

	"github.com/aneshas/tx"
	"github.com/aneshas/tx/example"
)

// NewAccount instantiates new Account postgres repository
func NewAccount(db *sql.DB) *Account {
	return &Account{
		SQL: tx.SQL{
			DB: db,
		},
	}
}

// Account represents account repo postrgres implementation
type Account struct {
	// Embed tx.SQL for a default sql tx.Transactional implementation
	// This step is optional, and you can implement tx.Transactional yourself
	// or embed tx.SQL and override any of the three methods (RunTx usually)
	tx.SQL
}

// ByID fetches new account for a given id
func (a *Account) ByID(ctx context.Context, id int64) (*example.Account, error) {
	// You might want to refactor this to a method eg. a.getDB(ctx) that returns
	// regular sql.DB if not in a tx or sql.Tx otherwise
	// (you might need to create an interface) eg: https://github.com/golang/go/issues/14468
	if tx, ok := tx.Current(ctx); ok {
		// This means we are in a transaction
		sqlTx := tx.Unwrap().(*sql.Tx)

		// Use sqlTx here instead of a.DB
		_ = sqlTx
	}

	panic("not impl")
}

// Save saves account to the db
func (a *Account) Save(ctx context.Context, acc *example.Account) error {
	if tx, ok := tx.Current(ctx); ok {
		// This means we are in a transaction
		sqlTx := tx.Unwrap().(*sql.Tx)

		// Use sqlTx here instead of a.DB
		_ = sqlTx
	}

	panic("not impl")
}
