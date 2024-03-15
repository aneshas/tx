//go:build integration
// +build integration

package pgxtxv5_test

import (
	"context"
	"fmt"
	"github.com/aneshas/tx"
	"github.com/aneshas/tx/pgxtxv5"
	"github.com/aneshas/tx/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	t := new(testing.T)

	db, _ = testutil.SetupDB(t)

	m.Run()
}

func TestShould_Commit_Pgx_Transaction(t *testing.T) {
	name := "success_pgx"

	doPgx(t, tx.New(pgxtxv5.NewDBFromPool(db)), name, false)
	testutil.AssertSuccess(t, db, name)
}

func TestShould_Rollback_Pgx_Transaction(t *testing.T) {
	name := "failure_pgx"

	doPgx(t, tx.New(pgxtxv5.NewDBFromPool(db)), name, true)
	testutil.AssertFailure(t, db, name)
}

func doPgx(t *testing.T, transactor *tx.TX, name string, fail bool) {
	t.Helper()

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		ttx, _ := pgxtxv5.From(ctx)

		_, err := ttx.Exec(ctx, `insert into cats (name) values($1)`, name)
		if err != nil {
			return err
		}

		if fail {
			return fmt.Errorf("db error")
		}

		return err
	})

	if !fail {
		assert.NoError(t, err)
	}
}
