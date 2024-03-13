// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Transaction is an autogenerated mock type for the Transaction type
type Transaction struct {
	mock.Mock
}

type Transaction_Expecter struct {
	mock *mock.Mock
}

func (_m *Transaction) EXPECT() *Transaction_Expecter {
	return &Transaction_Expecter{mock: &_m.Mock}
}

// Commit provides a mock function with given fields: ctx
func (_m *Transaction) Commit(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Commit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transaction_Commit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Commit'
type Transaction_Commit_Call struct {
	*mock.Call
}

// Commit is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Transaction_Expecter) Commit(ctx interface{}) *Transaction_Commit_Call {
	return &Transaction_Commit_Call{Call: _e.mock.On("Commit", ctx)}
}

func (_c *Transaction_Commit_Call) Run(run func(ctx context.Context)) *Transaction_Commit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Transaction_Commit_Call) Return(_a0 error) *Transaction_Commit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Transaction_Commit_Call) RunAndReturn(run func(context.Context) error) *Transaction_Commit_Call {
	_c.Call.Return(run)
	return _c
}

// Rollback provides a mock function with given fields: ctx
func (_m *Transaction) Rollback(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Rollback")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transaction_Rollback_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Rollback'
type Transaction_Rollback_Call struct {
	*mock.Call
}

// Rollback is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Transaction_Expecter) Rollback(ctx interface{}) *Transaction_Rollback_Call {
	return &Transaction_Rollback_Call{Call: _e.mock.On("Rollback", ctx)}
}

func (_c *Transaction_Rollback_Call) Run(run func(ctx context.Context)) *Transaction_Rollback_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Transaction_Rollback_Call) Return(_a0 error) *Transaction_Rollback_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Transaction_Rollback_Call) RunAndReturn(run func(context.Context) error) *Transaction_Rollback_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransaction creates a new instance of Transaction. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransaction(t interface {
	mock.TestingT
	Cleanup(func())
}) *Transaction {
	mock := &Transaction{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}