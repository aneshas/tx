//go:build integration
// +build integration

package sqltx_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aneshas/tx"
	"github.com/aneshas/tx/sqltx"
	"github.com/aneshas/tx/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	pool *pgxpool.Pool
	db   *sql.DB
)

func TestMain(m *testing.M) {
	t := new(testing.T)

	pool, db = testutil.SetupDB(t)

	m.Run()
}

func TestShould_Commit_Sql_Transaction(t *testing.T) {
	name := "success_sql"

	doSql(t, tx.New(sqltx.NewDB(db)), name, false)
	testutil.AssertSuccess(t, pool, name)
}

func TestShould_Rollback_Sql_Transaction(t *testing.T) {
	name := "failure_sql"

	doSql(t, tx.New(sqltx.NewDB(db)), name, true)
	testutil.AssertFailure(t, pool, name)
}

func doSql(t *testing.T, transactor *tx.TX, name string, fail bool) {
	t.Helper()

	err := transactor.Do(context.TODO(), func(ctx context.Context) error {
		ttx, _ := sqltx.From(ctx)

		_, err := ttx.Exec(`insert into cats (name) values($1)`, name)
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
