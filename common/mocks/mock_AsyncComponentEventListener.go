// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	common "github.com/couchbase/goxdcr/common"
	mock "github.com/stretchr/testify/mock"
)

// AsyncComponentEventListener is an autogenerated mock type for the AsyncComponentEventListener type
type AsyncComponentEventListener struct {
	mock.Mock
}

type AsyncComponentEventListener_Expecter struct {
	mock *mock.Mock
}

func (_m *AsyncComponentEventListener) EXPECT() *AsyncComponentEventListener_Expecter {
	return &AsyncComponentEventListener_Expecter{mock: &_m.Mock}
}

// Id provides a mock function with given fields:
func (_m *AsyncComponentEventListener) Id() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// AsyncComponentEventListener_Id_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Id'
type AsyncComponentEventListener_Id_Call struct {
	*mock.Call
}

// Id is a helper method to define mock.On call
func (_e *AsyncComponentEventListener_Expecter) Id() *AsyncComponentEventListener_Id_Call {
	return &AsyncComponentEventListener_Id_Call{Call: _e.mock.On("Id")}
}

func (_c *AsyncComponentEventListener_Id_Call) Run(run func()) *AsyncComponentEventListener_Id_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AsyncComponentEventListener_Id_Call) Return(_a0 string) *AsyncComponentEventListener_Id_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncComponentEventListener_Id_Call) RunAndReturn(run func() string) *AsyncComponentEventListener_Id_Call {
	_c.Call.Return(run)
	return _c
}

// OnEvent provides a mock function with given fields: event
func (_m *AsyncComponentEventListener) OnEvent(event *common.Event) {
	_m.Called(event)
}

// AsyncComponentEventListener_OnEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'OnEvent'
type AsyncComponentEventListener_OnEvent_Call struct {
	*mock.Call
}

// OnEvent is a helper method to define mock.On call
//   - event *common.Event
func (_e *AsyncComponentEventListener_Expecter) OnEvent(event interface{}) *AsyncComponentEventListener_OnEvent_Call {
	return &AsyncComponentEventListener_OnEvent_Call{Call: _e.mock.On("OnEvent", event)}
}

func (_c *AsyncComponentEventListener_OnEvent_Call) Run(run func(event *common.Event)) *AsyncComponentEventListener_OnEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*common.Event))
	})
	return _c
}

func (_c *AsyncComponentEventListener_OnEvent_Call) Return() *AsyncComponentEventListener_OnEvent_Call {
	_c.Call.Return()
	return _c
}

func (_c *AsyncComponentEventListener_OnEvent_Call) RunAndReturn(run func(*common.Event)) *AsyncComponentEventListener_OnEvent_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterComponentEventHandler provides a mock function with given fields: handler
func (_m *AsyncComponentEventListener) RegisterComponentEventHandler(handler common.AsyncComponentEventHandler) {
	_m.Called(handler)
}

// AsyncComponentEventListener_RegisterComponentEventHandler_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterComponentEventHandler'
type AsyncComponentEventListener_RegisterComponentEventHandler_Call struct {
	*mock.Call
}

// RegisterComponentEventHandler is a helper method to define mock.On call
//   - handler common.AsyncComponentEventHandler
func (_e *AsyncComponentEventListener_Expecter) RegisterComponentEventHandler(handler interface{}) *AsyncComponentEventListener_RegisterComponentEventHandler_Call {
	return &AsyncComponentEventListener_RegisterComponentEventHandler_Call{Call: _e.mock.On("RegisterComponentEventHandler", handler)}
}

func (_c *AsyncComponentEventListener_RegisterComponentEventHandler_Call) Run(run func(handler common.AsyncComponentEventHandler)) *AsyncComponentEventListener_RegisterComponentEventHandler_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.AsyncComponentEventHandler))
	})
	return _c
}

func (_c *AsyncComponentEventListener_RegisterComponentEventHandler_Call) Return() *AsyncComponentEventListener_RegisterComponentEventHandler_Call {
	_c.Call.Return()
	return _c
}

func (_c *AsyncComponentEventListener_RegisterComponentEventHandler_Call) RunAndReturn(run func(common.AsyncComponentEventHandler)) *AsyncComponentEventListener_RegisterComponentEventHandler_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields:
func (_m *AsyncComponentEventListener) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AsyncComponentEventListener_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type AsyncComponentEventListener_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
func (_e *AsyncComponentEventListener_Expecter) Start() *AsyncComponentEventListener_Start_Call {
	return &AsyncComponentEventListener_Start_Call{Call: _e.mock.On("Start")}
}

func (_c *AsyncComponentEventListener_Start_Call) Run(run func()) *AsyncComponentEventListener_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AsyncComponentEventListener_Start_Call) Return(_a0 error) *AsyncComponentEventListener_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncComponentEventListener_Start_Call) RunAndReturn(run func() error) *AsyncComponentEventListener_Start_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *AsyncComponentEventListener) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AsyncComponentEventListener_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type AsyncComponentEventListener_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *AsyncComponentEventListener_Expecter) Stop() *AsyncComponentEventListener_Stop_Call {
	return &AsyncComponentEventListener_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *AsyncComponentEventListener_Stop_Call) Run(run func()) *AsyncComponentEventListener_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AsyncComponentEventListener_Stop_Call) Return(_a0 error) *AsyncComponentEventListener_Stop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AsyncComponentEventListener_Stop_Call) RunAndReturn(run func() error) *AsyncComponentEventListener_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// NewAsyncComponentEventListener creates a new instance of AsyncComponentEventListener. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAsyncComponentEventListener(t interface {
	mock.TestingT
	Cleanup(func())
}) *AsyncComponentEventListener {
	mock := &AsyncComponentEventListener{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
