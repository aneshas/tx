package tx_test

import (
	"context"
	"fmt"
	"github.com/aneshas/tx/v2"
	"github.com/aneshas/tx/v2/testutil"
	"github.com/aneshas/tx/v2/testutil/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShould_Report_Transaction_Begin_Error(t *testing.T) {
	wantErr := fmt.Errorf("something bad occurred")

	db := testutil.NewDB(
		t,
		testutil.WithUnsuccessfulTransactionStart(wantErr),
	)
	transactor := tx.New(db)

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		return nil
	})

	assert.ErrorIs(t, err, wantErr)
}

func TestShould_Commit_Transaction_On_No_Error(t *testing.T) {
	db := testutil.NewDB(
		t,
		testutil.WithSuccessfulCommit(),
	)
	transactor := tx.New(db)

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		return nil
	})

	assert.NoError(t, err)
}

func TestShould_Rollback_Transaction_On_Error(t *testing.T) {
	db := testutil.NewDB(
		t,
		testutil.WithSuccessfulRollback(),
	)
	transactor := tx.New(db)

	wantErr := fmt.Errorf("something bad occurred")

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		return wantErr
	})

	assert.ErrorIs(t, err, wantErr)
}

func TestShould_Report_Unsuccessful_Rollback(t *testing.T) {
	wantTxErr := fmt.Errorf("something bad occurred")

	db := testutil.NewDB(
		t,
		testutil.WithUnsuccessfulRollback(wantTxErr),
	)
	transactor := tx.New(db)

	wantErr := fmt.Errorf("process error")

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		return wantErr
	})

	assert.ErrorIs(t, err, wantErr)
	assert.ErrorIs(t, err, wantTxErr)
}

func TestShould_Rollback_Transaction_On_Panic_And_RePanic(t *testing.T) {
	db := testutil.NewDB(
		t,
		testutil.WithSuccessfulRollback(),
	)
	transactor := tx.New(db)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic to be propagated")
		}
	}()

	_ = transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		panic("something very bad occurred")
	})
}

func TestShould_Still_Commit_On_Ignored_Error_And_Propagate_Error(t *testing.T) {
	wantErr := fmt.Errorf("something bad occurred")

	db := testutil.NewDB(
		t,
		testutil.WithSuccessfulCommit(),
	)
	transactor := tx.New(db, tx.WithIgnoredErrors(wantErr))

	err := transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		return wantErr
	})

	assert.ErrorIs(t, err, wantErr)
}

func TestShould_Retrieve_Tx_From_Context(t *testing.T) {
	db := testutil.NewDB(
		t,
		testutil.WithSuccessfulCommit(),
	)
	transactor := tx.New(db)

	_ = transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		ttx, ok := tx.From[tx.Transaction](ctx)

		assert.True(t, ok)
		assert.IsType(t, &mocks.Transaction{}, ttx)

		return nil
	})
}

func TestShould_Not_Retrieve_Conn_From_Context_On_Mismatched_Type(t *testing.T) {
	db := testutil.NewDB(
		t,
		testutil.WithSuccessfulCommit(),
	)
	transactor := tx.New(db)

	_ = transactor.WithTransaction(context.TODO(), func(ctx context.Context) error {
		_, ok := tx.From[mocks.Transaction](ctx)

		assert.False(t, ok)

		return nil
	})
}

func TestShould_Not_Retrieve_Conn_From_Context_Without_Transaction(t *testing.T) {
	ttx, ok := tx.From[*mocks.Transaction](context.TODO())

	assert.False(t, ok)
	assert.Nil(t, ttx)
}
