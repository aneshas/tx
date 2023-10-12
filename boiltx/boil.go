package boiltx

import (
	"context"
	"database/sql"
	"github.com/aneshas/tx"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// DB is sqlboiler helper function which returns the transaction stored in ctx as boil.ContextExecutor,
// otherwise the provided db is piped back
func DB(ctx context.Context, db *sql.DB) boil.ContextExecutor {
	sqlTx := tx.Tx(ctx)
	if sqlTx == nil {
		return db
	}

	return sqlTx
}
