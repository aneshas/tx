// Package tx provides a simple implementation agnostic transaction abstraction
package tx

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type ctxkey string

const txkey ctxkey = "_tonto_tx"

// Transactional interface should be implemented by an object wanting to use transactions
type Transactional interface {
	// RunTx should create new client side transaction object that can
	// be used to Commit or Rollback the transaction at a later point, and
	// wrap into *Tx by using tx.Wrap and then it should call tx.Run
	//
	// Example body:
	// {
	// 	t, _ := sql.BeginTx(ctx) // create client side transaction object
	// 	transaction := tx.Wrap(t) // wrap it inside of a tx transaction

	// TODO - Add any meta data you might need later to the transaction

	// 	return tx.Run(ctx, ur, transaction, f)
	// }
	RunTx(context.Context, func(context.Context) error) error

	// Commit is called every time RunTx transaction func returns nil error
	// indicating that transaction shuld be commited
	//
	// If Commit returns an error it is propagated and returned by RunTx
	//
	// Commit should never be called directly by client code
	Commit(*Tx) error

	// Rollback is called every time RunTx transaction func
	// returns an error indicating that transaction should be aborted
	//
	// If Rollback returns an error it wraps the error returned from
	// transaction func and both are propagated and returned by RunT
	//
	// Rollback should never be called directly by client code
	Rollback(*Tx) error
}

// Wrap creates new transaction and wraps provided client transaction
// to be used by the subsequent client side calls and ultimately
// by Commit and Rollback methods
func Wrap(clientTx interface{}) *Tx {
	return &Tx{
		clientTx: clientTx,
	}
}

// Tx represents transaction object
type Tx struct {
	m        sync.Mutex
	clientTx interface{}
}

// Unwrap unwraps the underlying client transaction wrapped by a call to tx.Wrap
func (tx *Tx) Unwrap() interface{} {
	tx.m.Lock()
	defer tx.m.Unlock()
	return tx.clientTx
}

// Current extracts Tx object from context if any
// Returns nil, false if not in a transaction
func Current(ctx context.Context) (*Tx, bool) {
	tx, ok := ctx.Value(txkey).(*Tx)
	return tx, ok
}

// Run runs provided func as a transaction
//
// Non nil error returned by the function is considered to indicate
// that the transaction should be rolled back, accordingly, nil error
// indicates that the transaction should be commited by calling
// Commit and Rollback methods on a registered repository respectively
func Run(ctx context.Context, t Transactional, tx *Tx, f func(context.Context) error) (err error) {

	// TODO - Add tx options
	// for _, opt := range opts {
	// 	opt(tx)
	// }

	defer func() {
		if err != nil {
			if e := t.Rollback(tx); e != nil {
				err = errors.Wrap(
					err,
					fmt.Sprintf("tx: error rolling back transaction: %v", e),
				)
			}
			return
		}

		if e := t.Commit(tx); e != nil {
			err = errors.Wrap(e, "tx: error commiting transaction")
		}
	}()

	return f(context.WithValue(ctx, txkey, tx))
}
