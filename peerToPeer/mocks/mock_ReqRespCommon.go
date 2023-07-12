// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	peerToPeer "github.com/couchbase/goxdcr/peerToPeer"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// ReqRespCommon is an autogenerated mock type for the ReqRespCommon type
type ReqRespCommon struct {
	mock.Mock
}

type ReqRespCommon_Expecter struct {
	mock *mock.Mock
}

func (_m *ReqRespCommon) EXPECT() *ReqRespCommon_Expecter {
	return &ReqRespCommon_Expecter{mock: &_m.Mock}
}

// DeSerialize provides a mock function with given fields: _a0
func (_m *ReqRespCommon) DeSerialize(_a0 []byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReqRespCommon_DeSerialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeSerialize'
type ReqRespCommon_DeSerialize_Call struct {
	*mock.Call
}

// DeSerialize is a helper method to define mock.On call
//   - _a0 []byte
func (_e *ReqRespCommon_Expecter) DeSerialize(_a0 interface{}) *ReqRespCommon_DeSerialize_Call {
	return &ReqRespCommon_DeSerialize_Call{Call: _e.mock.On("DeSerialize", _a0)}
}

func (_c *ReqRespCommon_DeSerialize_Call) Run(run func(_a0 []byte)) *ReqRespCommon_DeSerialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *ReqRespCommon_DeSerialize_Call) Return(_a0 error) *ReqRespCommon_DeSerialize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReqRespCommon_DeSerialize_Call) RunAndReturn(run func([]byte) error) *ReqRespCommon_DeSerialize_Call {
	_c.Call.Return(run)
	return _c
}

// GetEnqueuedTime provides a mock function with given fields:
func (_m *ReqRespCommon) GetEnqueuedTime() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// ReqRespCommon_GetEnqueuedTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEnqueuedTime'
type ReqRespCommon_GetEnqueuedTime_Call struct {
	*mock.Call
}

// GetEnqueuedTime is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) GetEnqueuedTime() *ReqRespCommon_GetEnqueuedTime_Call {
	return &ReqRespCommon_GetEnqueuedTime_Call{Call: _e.mock.On("GetEnqueuedTime")}
}

func (_c *ReqRespCommon_GetEnqueuedTime_Call) Run(run func()) *ReqRespCommon_GetEnqueuedTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_GetEnqueuedTime_Call) Return(_a0 time.Time) *ReqRespCommon_GetEnqueuedTime_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReqRespCommon_GetEnqueuedTime_Call) RunAndReturn(run func() time.Time) *ReqRespCommon_GetEnqueuedTime_Call {
	_c.Call.Return(run)
	return _c
}

// GetOpaque provides a mock function with given fields:
func (_m *ReqRespCommon) GetOpaque() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// ReqRespCommon_GetOpaque_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOpaque'
type ReqRespCommon_GetOpaque_Call struct {
	*mock.Call
}

// GetOpaque is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) GetOpaque() *ReqRespCommon_GetOpaque_Call {
	return &ReqRespCommon_GetOpaque_Call{Call: _e.mock.On("GetOpaque")}
}

func (_c *ReqRespCommon_GetOpaque_Call) Run(run func()) *ReqRespCommon_GetOpaque_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_GetOpaque_Call) Return(_a0 uint32) *ReqRespCommon_GetOpaque_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReqRespCommon_GetOpaque_Call) RunAndReturn(run func() uint32) *ReqRespCommon_GetOpaque_Call {
	_c.Call.Return(run)
	return _c
}

// GetOpcode provides a mock function with given fields:
func (_m *ReqRespCommon) GetOpcode() peerToPeer.OpCode {
	ret := _m.Called()

	var r0 peerToPeer.OpCode
	if rf, ok := ret.Get(0).(func() peerToPeer.OpCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(peerToPeer.OpCode)
	}

	return r0
}

// ReqRespCommon_GetOpcode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOpcode'
type ReqRespCommon_GetOpcode_Call struct {
	*mock.Call
}

// GetOpcode is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) GetOpcode() *ReqRespCommon_GetOpcode_Call {
	return &ReqRespCommon_GetOpcode_Call{Call: _e.mock.On("GetOpcode")}
}

func (_c *ReqRespCommon_GetOpcode_Call) Run(run func()) *ReqRespCommon_GetOpcode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_GetOpcode_Call) Return(_a0 peerToPeer.OpCode) *ReqRespCommon_GetOpcode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReqRespCommon_GetOpcode_Call) RunAndReturn(run func() peerToPeer.OpCode) *ReqRespCommon_GetOpcode_Call {
	_c.Call.Return(run)
	return _c
}

// GetSender provides a mock function with given fields:
func (_m *ReqRespCommon) GetSender() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ReqRespCommon_GetSender_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSender'
type ReqRespCommon_GetSender_Call struct {
	*mock.Call
}

// GetSender is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) GetSender() *ReqRespCommon_GetSender_Call {
	return &ReqRespCommon_GetSender_Call{Call: _e.mock.On("GetSender")}
}

func (_c *ReqRespCommon_GetSender_Call) Run(run func()) *ReqRespCommon_GetSender_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_GetSender_Call) Return(_a0 string) *ReqRespCommon_GetSender_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReqRespCommon_GetSender_Call) RunAndReturn(run func() string) *ReqRespCommon_GetSender_Call {
	_c.Call.Return(run)
	return _c
}

// GetType provides a mock function with given fields:
func (_m *ReqRespCommon) GetType() peerToPeer.ReqRespType {
	ret := _m.Called()

	var r0 peerToPeer.ReqRespType
	if rf, ok := ret.Get(0).(func() peerToPeer.ReqRespType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(peerToPeer.ReqRespType)
	}

	return r0
}

// ReqRespCommon_GetType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetType'
type ReqRespCommon_GetType_Call struct {
	*mock.Call
}

// GetType is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) GetType() *ReqRespCommon_GetType_Call {
	return &ReqRespCommon_GetType_Call{Call: _e.mock.On("GetType")}
}

func (_c *ReqRespCommon_GetType_Call) Run(run func()) *ReqRespCommon_GetType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_GetType_Call) Return(_a0 peerToPeer.ReqRespType) *ReqRespCommon_GetType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ReqRespCommon_GetType_Call) RunAndReturn(run func() peerToPeer.ReqRespType) *ReqRespCommon_GetType_Call {
	_c.Call.Return(run)
	return _c
}

// RecordEnqueuedTime provides a mock function with given fields:
func (_m *ReqRespCommon) RecordEnqueuedTime() {
	_m.Called()
}

// ReqRespCommon_RecordEnqueuedTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecordEnqueuedTime'
type ReqRespCommon_RecordEnqueuedTime_Call struct {
	*mock.Call
}

// RecordEnqueuedTime is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) RecordEnqueuedTime() *ReqRespCommon_RecordEnqueuedTime_Call {
	return &ReqRespCommon_RecordEnqueuedTime_Call{Call: _e.mock.On("RecordEnqueuedTime")}
}

func (_c *ReqRespCommon_RecordEnqueuedTime_Call) Run(run func()) *ReqRespCommon_RecordEnqueuedTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_RecordEnqueuedTime_Call) Return() *ReqRespCommon_RecordEnqueuedTime_Call {
	_c.Call.Return()
	return _c
}

func (_c *ReqRespCommon_RecordEnqueuedTime_Call) RunAndReturn(run func()) *ReqRespCommon_RecordEnqueuedTime_Call {
	_c.Call.Return(run)
	return _c
}

// Serialize provides a mock function with given fields:
func (_m *ReqRespCommon) Serialize() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]byte, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReqRespCommon_Serialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serialize'
type ReqRespCommon_Serialize_Call struct {
	*mock.Call
}

// Serialize is a helper method to define mock.On call
func (_e *ReqRespCommon_Expecter) Serialize() *ReqRespCommon_Serialize_Call {
	return &ReqRespCommon_Serialize_Call{Call: _e.mock.On("Serialize")}
}

func (_c *ReqRespCommon_Serialize_Call) Run(run func()) *ReqRespCommon_Serialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ReqRespCommon_Serialize_Call) Return(_a0 []byte, _a1 error) *ReqRespCommon_Serialize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ReqRespCommon_Serialize_Call) RunAndReturn(run func() ([]byte, error)) *ReqRespCommon_Serialize_Call {
	_c.Call.Return(run)
	return _c
}

// NewReqRespCommon creates a new instance of ReqRespCommon. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReqRespCommon(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReqRespCommon {
	mock := &ReqRespCommon{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
