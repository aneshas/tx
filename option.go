package tx

type Option func(tx *TX)

// WithIgnoredErrors offers a way to provide a list of errors which will
// not cause the transaction to be rolled back.
//
// The transaction will still be committed but the actual error will be returned
// by the WithTransaction method.
func WithIgnoredErrors(errs ...error) Option {
	return func(tx *TX) {
		tx.ignoreErrs = append(tx.ignoreErrs, errs...)
	}
}
