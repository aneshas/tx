package tx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aneshas/tx"
	"github.com/pkg/errors"
)

func TestCommit(t *testing.T) {
	cases := []struct {
		name    string
		t       *repoMock
		f       func(context.Context) error
		wantErr error
	}{
		{
			name: "commit success",
			t: &repoMock{
				CommitFunc: func(tx *tx.Tx) error {
					return nil
				},
			},
			f: func(ctx context.Context) error {
				return nil
			},
			wantErr: nil,
		},
		{
			name: "commit error",
			t: &repoMock{
				CommitFunc: func(tx *tx.Tx) error {
					return fmt.Errorf("db error")
				},
			},
			f: func(ctx context.Context) error {
				return nil
			},
			wantErr: errors.Wrap(
				fmt.Errorf("db error"),
				"tx: error commiting transaction",
			),
		},
		// Test context propagation
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			err := tc.t.RunTx(ctx, tc.f)
			if err != nil && tc.wantErr.Error() != err.Error() {
				t.Errorf("unexpected error response: want: %v got: %v", tc.wantErr, err)
			}
		})
	}
}

func TestRollback(t *testing.T) {

}

type repoMock struct {
	CommitFunc   func(*tx.Tx) error
	RollbackFunc func(*tx.Tx) error
}

func (rm *repoMock) RunTx(ctx context.Context, f func(context.Context) error) error {
	t := 64
	transaction := tx.Wrap(t)
	return tx.Run(ctx, rm, transaction, f)
}

func (rm *repoMock) Commit(tx *tx.Tx) error {
	return rm.CommitFunc(tx)
}

func (rm *repoMock) Rollback(tx *tx.Tx) error {
	return rm.RollbackFunc(tx)
}
