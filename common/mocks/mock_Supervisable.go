// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// Supervisable is an autogenerated mock type for the Supervisable type
type Supervisable struct {
	mock.Mock
}

type Supervisable_Expecter struct {
	mock *mock.Mock
}

func (_m *Supervisable) EXPECT() *Supervisable_Expecter {
	return &Supervisable_Expecter{mock: &_m.Mock}
}

// HeartBeat_async provides a mock function with given fields: respchan, timestamp
func (_m *Supervisable) HeartBeat_async(respchan chan []interface{}, timestamp time.Time) error {
	ret := _m.Called(respchan, timestamp)

	var r0 error
	if rf, ok := ret.Get(0).(func(chan []interface{}, time.Time) error); ok {
		r0 = rf(respchan, timestamp)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Supervisable_HeartBeat_async_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HeartBeat_async'
type Supervisable_HeartBeat_async_Call struct {
	*mock.Call
}

// HeartBeat_async is a helper method to define mock.On call
//   - respchan chan []interface{}
//   - timestamp time.Time
func (_e *Supervisable_Expecter) HeartBeat_async(respchan interface{}, timestamp interface{}) *Supervisable_HeartBeat_async_Call {
	return &Supervisable_HeartBeat_async_Call{Call: _e.mock.On("HeartBeat_async", respchan, timestamp)}
}

func (_c *Supervisable_HeartBeat_async_Call) Run(run func(respchan chan []interface{}, timestamp time.Time)) *Supervisable_HeartBeat_async_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(chan []interface{}), args[1].(time.Time))
	})
	return _c
}

func (_c *Supervisable_HeartBeat_async_Call) Return(_a0 error) *Supervisable_HeartBeat_async_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Supervisable_HeartBeat_async_Call) RunAndReturn(run func(chan []interface{}, time.Time) error) *Supervisable_HeartBeat_async_Call {
	_c.Call.Return(run)
	return _c
}

// HeartBeat_sync provides a mock function with given fields:
func (_m *Supervisable) HeartBeat_sync() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Supervisable_HeartBeat_sync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HeartBeat_sync'
type Supervisable_HeartBeat_sync_Call struct {
	*mock.Call
}

// HeartBeat_sync is a helper method to define mock.On call
func (_e *Supervisable_Expecter) HeartBeat_sync() *Supervisable_HeartBeat_sync_Call {
	return &Supervisable_HeartBeat_sync_Call{Call: _e.mock.On("HeartBeat_sync")}
}

func (_c *Supervisable_HeartBeat_sync_Call) Run(run func()) *Supervisable_HeartBeat_sync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Supervisable_HeartBeat_sync_Call) Return(_a0 bool) *Supervisable_HeartBeat_sync_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Supervisable_HeartBeat_sync_Call) RunAndReturn(run func() bool) *Supervisable_HeartBeat_sync_Call {
	_c.Call.Return(run)
	return _c
}

// Id provides a mock function with given fields:
func (_m *Supervisable) Id() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Supervisable_Id_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Id'
type Supervisable_Id_Call struct {
	*mock.Call
}

// Id is a helper method to define mock.On call
func (_e *Supervisable_Expecter) Id() *Supervisable_Id_Call {
	return &Supervisable_Id_Call{Call: _e.mock.On("Id")}
}

func (_c *Supervisable_Id_Call) Run(run func()) *Supervisable_Id_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Supervisable_Id_Call) Return(_a0 string) *Supervisable_Id_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Supervisable_Id_Call) RunAndReturn(run func() string) *Supervisable_Id_Call {
	_c.Call.Return(run)
	return _c
}

// IsReadyForHeartBeat provides a mock function with given fields:
func (_m *Supervisable) IsReadyForHeartBeat() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Supervisable_IsReadyForHeartBeat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsReadyForHeartBeat'
type Supervisable_IsReadyForHeartBeat_Call struct {
	*mock.Call
}

// IsReadyForHeartBeat is a helper method to define mock.On call
func (_e *Supervisable_Expecter) IsReadyForHeartBeat() *Supervisable_IsReadyForHeartBeat_Call {
	return &Supervisable_IsReadyForHeartBeat_Call{Call: _e.mock.On("IsReadyForHeartBeat")}
}

func (_c *Supervisable_IsReadyForHeartBeat_Call) Run(run func()) *Supervisable_IsReadyForHeartBeat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Supervisable_IsReadyForHeartBeat_Call) Return(_a0 bool) *Supervisable_IsReadyForHeartBeat_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Supervisable_IsReadyForHeartBeat_Call) RunAndReturn(run func() bool) *Supervisable_IsReadyForHeartBeat_Call {
	_c.Call.Return(run)
	return _c
}

// NewSupervisable creates a new instance of Supervisable. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSupervisable(t interface {
	mock.TestingT
	Cleanup(func())
}) *Supervisable {
	mock := &Supervisable{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
