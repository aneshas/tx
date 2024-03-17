package testtx_test

import (
	"context"
	"fmt"
	"github.com/aneshas/tx/v2"
	"github.com/aneshas/tx/v2/testtx"
	"github.com/stretchr/testify/assert"
	"testing"
)

type svc struct {
	tx.Transactor

	DidSomething bool
}

func (s *svc) doSomething(ctx context.Context, err error) error {
	return s.WithTransaction(ctx, func(ctx context.Context) error {
		s.DidSomething = true

		return err
	})
}

func TestShould_Delegate_Call(t *testing.T) {
	transactor := testtx.New()

	s := &svc{
		Transactor: transactor,
	}

	err := s.doSomething(context.TODO(), nil)

	assert.NoError(t, err)
	assert.True(t, s.DidSomething)
}

func TestShould_Save_Err(t *testing.T) {
	transactor := testtx.New()

	s := &svc{
		Transactor: transactor,
	}

	wantErr := fmt.Errorf("something bad ocurred")

	err := s.doSomething(context.TODO(), wantErr)

	assert.ErrorIs(t, err, wantErr)
	assert.ErrorIs(t, transactor.Err, wantErr)
}
