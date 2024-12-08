package testtx

import "context"

// New creates a new TX
func New() *TX {
	return &TX{}
}

// TX is a noop test implementation of tx.DB
type TX struct {
	Err error
}

// WithTransaction is a noop test implementation of tx.DB
func (t *TX) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	t.Err = f(ctx)

	return t.Err
}
