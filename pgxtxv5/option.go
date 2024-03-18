package pgxtxv5

import "github.com/jackc/pgx/v5"

// PgxTxOption represents pgx driver transaction option
type PgxTxOption func(pool *Pool)

// WithTxOptions allows us to set transaction options (eg. isolation level)
func WithTxOptions(txOptions pgx.TxOptions) PgxTxOption {
	return func(pool *Pool) {
		pool.txOpts = txOptions
	}
}
