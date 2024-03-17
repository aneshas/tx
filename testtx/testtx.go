package testtx

import "context"

func New() *TX {
	return &TX{}
}

type TX struct {
	Err error
}

func (t *TX) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	t.Err = f(ctx)

	return t.Err
}
