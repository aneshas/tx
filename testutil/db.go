package testutil

import (
	"github.com/aneshas/tx/testutil/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func NewDB(t *testing.T, opts ...Option) *DB {
	db := DB{
		t:  t,
		DB: mocks.NewDB(t),
	}

	for _, opt := range opts {
		opt(&db)
	}

	return &db
}

type Option func(db *DB)

func WithUnsuccessfulTransactionStart(with error) Option {
	return func(db *DB) {
		db.EXPECT().Begin(mock.Anything).Return(nil, with).Once()
	}
}

func WithSuccessfulTransactionStart() Option {
	return func(db *DB) {
		tx := mocks.NewTransaction(db.t)
		db.EXPECT().Begin(mock.Anything).Return(tx, nil).Once()
	}
}

func WithSuccessfulCommit() Option {
	return func(db *DB) {
		tx := mocks.NewTransaction(db.t)
		tx.EXPECT().Commit(mock.Anything).Return(nil).Once()
		db.EXPECT().Begin(mock.Anything).Return(tx, nil).Once()
	}
}

func WithSuccessfulRollback() Option {
	return func(db *DB) {
		tx := mocks.NewTransaction(db.t)
		tx.EXPECT().Rollback(mock.Anything).Return(nil).Once()
		db.EXPECT().Begin(mock.Anything).Return(tx, nil).Once()
	}
}

func WithUnsuccessfulRollback(with error) Option {
	return func(db *DB) {
		tx := mocks.NewTransaction(db.t)
		tx.EXPECT().Rollback(mock.Anything).Return(with).Once()
		db.EXPECT().Begin(mock.Anything).Return(tx, nil).Once()
	}
}

type DB struct {
	t *testing.T
	*mocks.DB
}
