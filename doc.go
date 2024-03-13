// Package tx provides a simple transaction abstraction in order to enable decoupling/abstraction of persistence from
// application/domain logic while still leaving transaction control to the application service.
// (Something like @Transactional annotation in Java, without an annotation)
package tx

//go:generate mockery --all --with-expecter --case=underscore --output ./testutil/mocks
