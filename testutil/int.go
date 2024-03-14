package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

func SetupDB(t *testing.T) (*pgxpool.Pool, *sql.DB) {
	t.Helper()

	p := postgres.Preset(
		postgres.WithUser("gnomock", "gnomick"),
		postgres.WithDatabase("mydb"),
		postgres.WithQueriesFile("../testutil/db/schema.sql"),
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

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)

	return pool, db
}

func AssertSuccess(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()

	row := pool.QueryRow(context.TODO(), `select name from cats where name=$1`, name)

	var n string

	err := row.Scan(&n)
	assert.NoError(t, err)

	assert.Equal(t, name, n)
}

func AssertFailure(t *testing.T, pool *pgxpool.Pool, name string) {
	t.Helper()

	row := pool.QueryRow(context.TODO(), `select name from cats where name=$1`, name)

	var n string

	err := row.Scan(&n)

	assert.ErrorIs(t, err, pgx.ErrNoRows)
}
