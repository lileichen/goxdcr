// Code generated by mockery. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ThroughputThrottlerSvc is an autogenerated mock type for the ThroughputThrottlerSvc type
type ThroughputThrottlerSvc struct {
	mock.Mock
}

type ThroughputThrottlerSvc_Expecter struct {
	mock *mock.Mock
}

func (_m *ThroughputThrottlerSvc) EXPECT() *ThroughputThrottlerSvc_Expecter {
	return &ThroughputThrottlerSvc_Expecter{mock: &_m.Mock}
}

// CanSend provides a mock function with given fields: isHighPriorityReplication
func (_m *ThroughputThrottlerSvc) CanSend(isHighPriorityReplication bool) bool {
	ret := _m.Called(isHighPriorityReplication)

	var r0 bool
	if rf, ok := ret.Get(0).(func(bool) bool); ok {
		r0 = rf(isHighPriorityReplication)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// ThroughputThrottlerSvc_CanSend_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CanSend'
type ThroughputThrottlerSvc_CanSend_Call struct {
	*mock.Call
}

// CanSend is a helper method to define mock.On call
//   - isHighPriorityReplication bool
func (_e *ThroughputThrottlerSvc_Expecter) CanSend(isHighPriorityReplication interface{}) *ThroughputThrottlerSvc_CanSend_Call {
	return &ThroughputThrottlerSvc_CanSend_Call{Call: _e.mock.On("CanSend", isHighPriorityReplication)}
}

func (_c *ThroughputThrottlerSvc_CanSend_Call) Run(run func(isHighPriorityReplication bool)) *ThroughputThrottlerSvc_CanSend_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(bool))
	})
	return _c
}

func (_c *ThroughputThrottlerSvc_CanSend_Call) Return(_a0 bool) *ThroughputThrottlerSvc_CanSend_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThroughputThrottlerSvc_CanSend_Call) RunAndReturn(run func(bool) bool) *ThroughputThrottlerSvc_CanSend_Call {
	_c.Call.Return(run)
	return _c
}

// Start provides a mock function with given fields:
func (_m *ThroughputThrottlerSvc) Start() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ThroughputThrottlerSvc_Start_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Start'
type ThroughputThrottlerSvc_Start_Call struct {
	*mock.Call
}

// Start is a helper method to define mock.On call
func (_e *ThroughputThrottlerSvc_Expecter) Start() *ThroughputThrottlerSvc_Start_Call {
	return &ThroughputThrottlerSvc_Start_Call{Call: _e.mock.On("Start")}
}

func (_c *ThroughputThrottlerSvc_Start_Call) Run(run func()) *ThroughputThrottlerSvc_Start_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ThroughputThrottlerSvc_Start_Call) Return(_a0 error) *ThroughputThrottlerSvc_Start_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThroughputThrottlerSvc_Start_Call) RunAndReturn(run func() error) *ThroughputThrottlerSvc_Start_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *ThroughputThrottlerSvc) Stop() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ThroughputThrottlerSvc_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type ThroughputThrottlerSvc_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *ThroughputThrottlerSvc_Expecter) Stop() *ThroughputThrottlerSvc_Stop_Call {
	return &ThroughputThrottlerSvc_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *ThroughputThrottlerSvc_Stop_Call) Run(run func()) *ThroughputThrottlerSvc_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ThroughputThrottlerSvc_Stop_Call) Return(_a0 error) *ThroughputThrottlerSvc_Stop_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThroughputThrottlerSvc_Stop_Call) RunAndReturn(run func() error) *ThroughputThrottlerSvc_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateSettings provides a mock function with given fields: setting
func (_m *ThroughputThrottlerSvc) UpdateSettings(setting map[string]interface{}) map[string]error {
	ret := _m.Called(setting)

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func(map[string]interface{}) map[string]error); ok {
		r0 = rf(setting)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// ThroughputThrottlerSvc_UpdateSettings_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateSettings'
type ThroughputThrottlerSvc_UpdateSettings_Call struct {
	*mock.Call
}

// UpdateSettings is a helper method to define mock.On call
//   - setting map[string]interface{}
func (_e *ThroughputThrottlerSvc_Expecter) UpdateSettings(setting interface{}) *ThroughputThrottlerSvc_UpdateSettings_Call {
	return &ThroughputThrottlerSvc_UpdateSettings_Call{Call: _e.mock.On("UpdateSettings", setting)}
}

func (_c *ThroughputThrottlerSvc_UpdateSettings_Call) Run(run func(setting map[string]interface{})) *ThroughputThrottlerSvc_UpdateSettings_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(map[string]interface{}))
	})
	return _c
}

func (_c *ThroughputThrottlerSvc_UpdateSettings_Call) Return(_a0 map[string]error) *ThroughputThrottlerSvc_UpdateSettings_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ThroughputThrottlerSvc_UpdateSettings_Call) RunAndReturn(run func(map[string]interface{}) map[string]error) *ThroughputThrottlerSvc_UpdateSettings_Call {
	_c.Call.Return(run)
	return _c
}

// Wait provides a mock function with given fields:
func (_m *ThroughputThrottlerSvc) Wait() {
	_m.Called()
}

// ThroughputThrottlerSvc_Wait_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Wait'
type ThroughputThrottlerSvc_Wait_Call struct {
	*mock.Call
}

// Wait is a helper method to define mock.On call
func (_e *ThroughputThrottlerSvc_Expecter) Wait() *ThroughputThrottlerSvc_Wait_Call {
	return &ThroughputThrottlerSvc_Wait_Call{Call: _e.mock.On("Wait")}
}

func (_c *ThroughputThrottlerSvc_Wait_Call) Run(run func()) *ThroughputThrottlerSvc_Wait_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ThroughputThrottlerSvc_Wait_Call) Return() *ThroughputThrottlerSvc_Wait_Call {
	_c.Call.Return()
	return _c
}

func (_c *ThroughputThrottlerSvc_Wait_Call) RunAndReturn(run func()) *ThroughputThrottlerSvc_Wait_Call {
	_c.Call.Return(run)
	return _c
}

// NewThroughputThrottlerSvc creates a new instance of ThroughputThrottlerSvc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewThroughputThrottlerSvc(t interface {
	mock.TestingT
	Cleanup(func())
}) *ThroughputThrottlerSvc {
	mock := &ThroughputThrottlerSvc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
