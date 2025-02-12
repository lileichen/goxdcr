// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	mock "github.com/stretchr/testify/mock"
)

// CheckpointFunc is an autogenerated mock type for the CheckpointFunc type
type CheckpointFunc struct {
	mock.Mock
}

type CheckpointFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *CheckpointFunc) EXPECT() *CheckpointFunc_Expecter {
	return &CheckpointFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *CheckpointFunc) Execute(_a0 common.Pipeline) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Pipeline) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CheckpointFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type CheckpointFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 common.Pipeline
func (_e *CheckpointFunc_Expecter) Execute(_a0 interface{}) *CheckpointFunc_Execute_Call {
	return &CheckpointFunc_Execute_Call{Call: _e.mock.On("Execute", _a0)}
}

func (_c *CheckpointFunc_Execute_Call) Run(run func(_a0 common.Pipeline)) *CheckpointFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Pipeline))
	})
	return _c
}

func (_c *CheckpointFunc_Execute_Call) Return(_a0 error) *CheckpointFunc_Execute_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CheckpointFunc_Execute_Call) RunAndReturn(run func(common.Pipeline) error) *CheckpointFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewCheckpointFunc creates a new instance of CheckpointFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCheckpointFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *CheckpointFunc {
	mock := &CheckpointFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
