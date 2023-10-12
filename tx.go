// Package tx provides a simple transaction abstraction in order to enable decoupling/abstraction of persistence from
// application/domain logic while still leaving transaction control to the application service.
// (Something like @Transactional annotation in Java, without an annotation)
package tx

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type key struct{}

// Transactor is a helper transactor interface added for brevity purposes, so you don't have to define your own
type Transactor interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}

// New constructs new transactor which will use provided db to handle the transaction
func New(db *sql.DB) *TX {
	return &TX{db: db}
}

// TX represents sql transactor
type TX struct {
	db *sql.DB
}

// Do will execute func f in a sql transaction.
// This is mostly useful for when we want to control the transaction scope from
// application layer, for example application service/command handler.
// If f fails with an error, transactor will automatically try to roll back the transaction and report back any errors,
// otherwise the implicit transaction will be committed.
func (t *TX) Do(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil) // add tx options if different isolation levels are needed
	if err != nil {
		return fmt.Errorf("tx: could not start transaction: %w", err)
	}

	ctx = context.WithValue(ctx, key{}, tx)

	err = f(ctx)
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			return errors.Join(e, err)
		}

		return err
	}

	return tx.Commit()
}

// DB returns the underlying *sql.DB
func (t *TX) DB() *sql.DB {
	return t.db
}

// Tx returns the implicit transaction if available, otherwise it returns nil
func Tx(ctx context.Context) *sql.Tx {
	val := ctx.Value(key{})
	if val == nil {
		return nil
	}

	tx, ok := val.(*sql.Tx)
	if !ok {
		return nil
	}

	return tx
}
