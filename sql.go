package tx

import (
	"context"
	"database/sql"
)

// SQL is a sql implementation of Transactional interface
// It is meant to be used as a helper and embedded inside of data access objects such
// as a repositories in order to provide them with decoupled sql transaction behavior
type SQL struct {
	DB *sql.DB
}

// RunTx runs f in a SQL transaction
func (sql *SQL) RunTx(ctx context.Context, f func(context.Context) error) error {
	tx, err := sql.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	return Run(ctx, sql, Wrap(tx), f)
}

// Commit commits sql transaction
func (SQL) Commit(tx *Tx) error {
	transaction := tx.Unwrap().(*sql.Tx)
	return transaction.Commit()
}

// Rollback rolls back sql transaction
func (SQL) Rollback(tx *Tx) error {
	transaction := tx.Unwrap().(*sql.Tx)
	return transaction.Rollback()
}
