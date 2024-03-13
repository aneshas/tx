package tx

import (
	"context"
	"errors"
	"fmt"
)

type key struct{}

// Transactor is a helper transactor interface added for brevity purposes, so you don't have to define your own
type Transactor interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}

// DB represents an interface to a db capable of starting a transaction
type DB interface {
	Begin(ctx context.Context) (Transaction, error)
}

// Transaction represents db transaction
type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// New constructs new transactor which will use provided db to handle the transaction
func New(db DB, opts ...Option) *TX {
	ttx := TX{db: db}

	for _, opt := range opts {
		opt(&ttx)
	}

	return &ttx
}

// TX represents sql transactor
type TX struct {
	db         DB
	ignoreErrs []error
}

// Do will execute func f in a sql transaction.
// This is mostly useful for when we want to control the transaction scope from
// application layer, for example application service/command handler.
// If f fails with an error, transactor will automatically try to roll back the transaction and report back any errors,
// otherwise the implicit transaction will be committed.
func (t *TX) Do(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx) // add tx options if different isolation levels are needed
	if err != nil {
		return fmt.Errorf("tx: could not start transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	ctx = context.WithValue(ctx, key{}, tx)

	err = f(ctx)
	if err != nil && !t.shouldIgnore(err) {
		e := tx.Rollback(ctx)
		if e != nil {
			return errors.Join(e, err)
		}

		return err
	}

	return errors.Join(err, tx.Commit(ctx))
}

func (t *TX) shouldIgnore(err error) bool {
	for _, e := range t.ignoreErrs {
		if errors.Is(err, e) {
			return true
		}
	}

	return false
}

func From[T any](ctx context.Context, orDB DB) (T, error) {
	if t, ok := Conn[T](ctx); ok {
		return t.(T), nil
	}

	ttx, err := orDB.Begin(ctx)
	if err != nil {
		var t T

		return t, err
	}

	return ttx.(T), nil
}

// Conn returns underlying tx value from context if it can be type-casted to T
// Otherwise it returns nil, false
func Conn[T any](ctx context.Context) (any, bool) {
	val := ctx.Value(key{})
	if val == nil {
		return nil, false
	}

	tx, ok := val.(T)
	if !ok {
		return nil, false
	}

	return tx, true
}
