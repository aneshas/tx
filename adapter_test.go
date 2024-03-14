//go:build integration
// +build integration

package tx_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aneshas/tx"
	"github.com/aneshas/tx/pgxtxv5"
	"github.com/aneshas/tx/sqltx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
)

var (
	db    *pgxpool.Pool
	sqlDB *sql.DB
)

func TestMain(m *testing.M) {
	t := new(testing.T)

	setupDB(t)

	m.Run()
}

func setupDB(t *testing.T) {
	t.Helper()

	p := postgres.Preset(
		postgres.WithUser("gnomock", "gnomick"),
		postgres.WithDatabase("mydb"),
		postgres.WithQueriesFile("./testutil/db/schema.sql"),
	)

	container, err := gnomock.Start(p)
	assert.NoError(t, err)

	t.Cleanup(func() { _ = gnomock.Stop(container) })

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		container.Host, container.DefaultPort(),
		"gnomock", "gnomick", "mydb",
	)

	pgConfig, err := pgxpool.ParseConfig(connStr)
	assert.NoError(t, err)

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	assert.NoError(t, err)

	dbs, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)

	db = pool
	sqlDB = dbs
}

func TestShould_Commit_Pgx_Transaction(t *testing.T) {
	name := "success_pgx"

	doPgx(t, tx.New(pgxtxv5.NewDBFromPool(db)), db, name, false)
	assertSuccess(t, db, name)
}

func TestShould_Rollback_Pgx_Transaction(t *testing.T) {
	name := "failure_pgx"

	doPgx(t, tx.New(pgxtxv5.NewDBFromPool(db)), db, name, true)
	assertFailure(t, db, name)
}

func TestShould_Commit_Sql_Transaction(t *testing.T) {
	name := "success_sql"

	doSql(t, tx.New(sqltx.NewDB(sqlDB)), sqlDB, name, false)
	assertSuccess(t, db, name)
}

func TestShould_Rollback_Sql_Transaction(t *testing.T) {
	name := "failure_sql"

	doSql(t, tx.New(sqltx.NewDB(sqlDB)), sqlDB, name, true)
	assertFailure(t, db, name)
}

func doPgx(t *testing.T, transactor *tx.TX, pool *pgxpool.Pool, name string, fail bool) {
	t.Helper()

	err := transactor.Do(context.TODO(), func(ctx context.Context) error {
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

func doSql(t *testing.T, transactor *tx.TX, db *sql.DB, name string, fail bool) {
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

func assertSuccess(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()

	row := pool.QueryRow(context.TODO(), `select name from cats where name=$1`, name)

	var n string

	err := row.Scan(&n)
	assert.NoError(t, err)

	assert.Equal(t, name, n)
}

func assertFailure(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()

	row := pool.QueryRow(context.TODO(), `select name from cats where name=$1`, name)

	var n string

	err := row.Scan(&n)

	assert.ErrorIs(t, err, pgx.ErrNoRows)
}
