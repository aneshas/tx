package gormtx

import (
	"context"
	"github.com/aneshas/tx/v2"
	"gorm.io/gorm"
)

var (
	_ tx.DB          = &DB{}
	_ tx.Transaction = &Tx{}
)

// NewDB instantiates new tx.DB *gorm.DB wrapper
func NewDB(db *gorm.DB) tx.DB {
	return &DB{DB: db}
}

// DB implements tx.DB
type DB struct {
	*gorm.DB
}

// Begin begins gorm transaction
func (db *DB) Begin(ctx context.Context) (tx.Transaction, error) {
	txx := db.WithContext(ctx).Begin()
	if txx.Error != nil {
		return nil, txx.Error
	}

	return &Tx{txx}, nil
}

// Tx wraps *gorm.DB in order top implement tx.Transaction
type Tx struct {
	*gorm.DB
}

// Commit commits the transaction
func (t Tx) Commit(_ context.Context) error {
	return t.DB.Commit().Error
}

// Rollback rolls back the transaction
func (t Tx) Rollback(_ context.Context) error {
	return t.DB.Rollback().Error
}

// From returns underlying *gorm.DB (wrapped in *Tx)
func From(ctx context.Context) (*Tx, bool) {
	return tx.From[*Tx](ctx)
}
