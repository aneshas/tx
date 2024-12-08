//go:build integration
// +build integration

package gormtx_test

import (
	"context"
	"fmt"
	"github.com/aneshas/tx/v2"
	"github.com/aneshas/tx/v2/gormtx"
	"github.com/aneshas/tx/v2/testutil"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var (
	pool *pgxpool.Pool
	db   *gorm.DB
)

func TestMain(m *testing.M) {
	t := new(testing.T)

	p, sqlDB := testutil.SetupDB(t)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	assert.NoError(t, err)

	pool = p
	db = gormDB

	m.Run()
}

func TestShould_Commit_Sql_Transaction(t *testing.T) {
	name := "success_sql"

	doSql(t, tx.New(gormtx.NewDB(db)), name, false)
	testutil.AssertSuccess(t, pool, name)
}

func TestShould_Rollback_Sql_Transaction(t *testing.T) {
	name := "failure_sql"

	doSql(t, tx.New(gormtx.NewDB(db)), name, true)
	testutil.AssertFailure(t, pool, name)
}

func doSql(t *testing.T, transactor *tx.TX, name string, fail bool) {
	t.Helper()

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		ttx, _ := gormtx.From(ctx)

		db := ttx.Exec(`insert into cats (name) values(?)`, name)
		if db.Error != nil {
			return db.Error
		}

		if fail {
			return fmt.Errorf("db error")
		}

		return db.Error
	})

	if !fail {
		assert.NoError(t, err)
	}
}
