// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	peerToPeer "github.com/couchbase/goxdcr/peerToPeer"
	mock "github.com/stretchr/testify/mock"

	time "time"
)

// Response is an autogenerated mock type for the Response type
type Response struct {
	mock.Mock
}

type Response_Expecter struct {
	mock *mock.Mock
}

func (_m *Response) EXPECT() *Response_Expecter {
	return &Response_Expecter{mock: &_m.Mock}
}

// DeSerialize provides a mock function with given fields: _a0
func (_m *Response) DeSerialize(_a0 []byte) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Response_DeSerialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeSerialize'
type Response_DeSerialize_Call struct {
	*mock.Call
}

// DeSerialize is a helper method to define mock.On call
//   - _a0 []byte
func (_e *Response_Expecter) DeSerialize(_a0 interface{}) *Response_DeSerialize_Call {
	return &Response_DeSerialize_Call{Call: _e.mock.On("DeSerialize", _a0)}
}

func (_c *Response_DeSerialize_Call) Run(run func(_a0 []byte)) *Response_DeSerialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *Response_DeSerialize_Call) Return(_a0 error) *Response_DeSerialize_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_DeSerialize_Call) RunAndReturn(run func([]byte) error) *Response_DeSerialize_Call {
	_c.Call.Return(run)
	return _c
}

// GetEnqueuedTime provides a mock function with given fields:
func (_m *Response) GetEnqueuedTime() time.Time {
	ret := _m.Called()

	var r0 time.Time
	if rf, ok := ret.Get(0).(func() time.Time); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Time)
	}

	return r0
}

// Response_GetEnqueuedTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEnqueuedTime'
type Response_GetEnqueuedTime_Call struct {
	*mock.Call
}

// GetEnqueuedTime is a helper method to define mock.On call
func (_e *Response_Expecter) GetEnqueuedTime() *Response_GetEnqueuedTime_Call {
	return &Response_GetEnqueuedTime_Call{Call: _e.mock.On("GetEnqueuedTime")}
}

func (_c *Response_GetEnqueuedTime_Call) Run(run func()) *Response_GetEnqueuedTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_GetEnqueuedTime_Call) Return(_a0 time.Time) *Response_GetEnqueuedTime_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_GetEnqueuedTime_Call) RunAndReturn(run func() time.Time) *Response_GetEnqueuedTime_Call {
	_c.Call.Return(run)
	return _c
}

// GetErrorString provides a mock function with given fields:
func (_m *Response) GetErrorString() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Response_GetErrorString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetErrorString'
type Response_GetErrorString_Call struct {
	*mock.Call
}

// GetErrorString is a helper method to define mock.On call
func (_e *Response_Expecter) GetErrorString() *Response_GetErrorString_Call {
	return &Response_GetErrorString_Call{Call: _e.mock.On("GetErrorString")}
}

func (_c *Response_GetErrorString_Call) Run(run func()) *Response_GetErrorString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_GetErrorString_Call) Return(_a0 string) *Response_GetErrorString_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_GetErrorString_Call) RunAndReturn(run func() string) *Response_GetErrorString_Call {
	_c.Call.Return(run)
	return _c
}

// GetOpaque provides a mock function with given fields:
func (_m *Response) GetOpaque() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// Response_GetOpaque_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOpaque'
type Response_GetOpaque_Call struct {
	*mock.Call
}

// GetOpaque is a helper method to define mock.On call
func (_e *Response_Expecter) GetOpaque() *Response_GetOpaque_Call {
	return &Response_GetOpaque_Call{Call: _e.mock.On("GetOpaque")}
}

func (_c *Response_GetOpaque_Call) Run(run func()) *Response_GetOpaque_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_GetOpaque_Call) Return(_a0 uint32) *Response_GetOpaque_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_GetOpaque_Call) RunAndReturn(run func() uint32) *Response_GetOpaque_Call {
	_c.Call.Return(run)
	return _c
}

// GetOpcode provides a mock function with given fields:
func (_m *Response) GetOpcode() peerToPeer.OpCode {
	ret := _m.Called()

	var r0 peerToPeer.OpCode
	if rf, ok := ret.Get(0).(func() peerToPeer.OpCode); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(peerToPeer.OpCode)
	}

	return r0
}

// Response_GetOpcode_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOpcode'
type Response_GetOpcode_Call struct {
	*mock.Call
}

// GetOpcode is a helper method to define mock.On call
func (_e *Response_Expecter) GetOpcode() *Response_GetOpcode_Call {
	return &Response_GetOpcode_Call{Call: _e.mock.On("GetOpcode")}
}

func (_c *Response_GetOpcode_Call) Run(run func()) *Response_GetOpcode_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_GetOpcode_Call) Return(_a0 peerToPeer.OpCode) *Response_GetOpcode_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_GetOpcode_Call) RunAndReturn(run func() peerToPeer.OpCode) *Response_GetOpcode_Call {
	_c.Call.Return(run)
	return _c
}

// GetSender provides a mock function with given fields:
func (_m *Response) GetSender() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Response_GetSender_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSender'
type Response_GetSender_Call struct {
	*mock.Call
}

// GetSender is a helper method to define mock.On call
func (_e *Response_Expecter) GetSender() *Response_GetSender_Call {
	return &Response_GetSender_Call{Call: _e.mock.On("GetSender")}
}

func (_c *Response_GetSender_Call) Run(run func()) *Response_GetSender_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_GetSender_Call) Return(_a0 string) *Response_GetSender_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_GetSender_Call) RunAndReturn(run func() string) *Response_GetSender_Call {
	_c.Call.Return(run)
	return _c
}

// GetType provides a mock function with given fields:
func (_m *Response) GetType() peerToPeer.ReqRespType {
	ret := _m.Called()

	var r0 peerToPeer.ReqRespType
	if rf, ok := ret.Get(0).(func() peerToPeer.ReqRespType); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(peerToPeer.ReqRespType)
	}

	return r0
}

// Response_GetType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetType'
type Response_GetType_Call struct {
	*mock.Call
}

// GetType is a helper method to define mock.On call
func (_e *Response_Expecter) GetType() *Response_GetType_Call {
	return &Response_GetType_Call{Call: _e.mock.On("GetType")}
}

func (_c *Response_GetType_Call) Run(run func()) *Response_GetType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_GetType_Call) Return(_a0 peerToPeer.ReqRespType) *Response_GetType_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Response_GetType_Call) RunAndReturn(run func() peerToPeer.ReqRespType) *Response_GetType_Call {
	_c.Call.Return(run)
	return _c
}

// RecordEnqueuedTime provides a mock function with given fields:
func (_m *Response) RecordEnqueuedTime() {
	_m.Called()
}

// Response_RecordEnqueuedTime_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RecordEnqueuedTime'
type Response_RecordEnqueuedTime_Call struct {
	*mock.Call
}

// RecordEnqueuedTime is a helper method to define mock.On call
func (_e *Response_Expecter) RecordEnqueuedTime() *Response_RecordEnqueuedTime_Call {
	return &Response_RecordEnqueuedTime_Call{Call: _e.mock.On("RecordEnqueuedTime")}
}

func (_c *Response_RecordEnqueuedTime_Call) Run(run func()) *Response_RecordEnqueuedTime_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_RecordEnqueuedTime_Call) Return() *Response_RecordEnqueuedTime_Call {
	_c.Call.Return()
	return _c
}

func (_c *Response_RecordEnqueuedTime_Call) RunAndReturn(run func()) *Response_RecordEnqueuedTime_Call {
	_c.Call.Return(run)
	return _c
}

// Serialize provides a mock function with given fields:
func (_m *Response) Serialize() ([]byte, error) {
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

// Response_Serialize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Serialize'
type Response_Serialize_Call struct {
	*mock.Call
}

// Serialize is a helper method to define mock.On call
func (_e *Response_Expecter) Serialize() *Response_Serialize_Call {
	return &Response_Serialize_Call{Call: _e.mock.On("Serialize")}
}

func (_c *Response_Serialize_Call) Run(run func()) *Response_Serialize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Response_Serialize_Call) Return(_a0 []byte, _a1 error) *Response_Serialize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Response_Serialize_Call) RunAndReturn(run func() ([]byte, error)) *Response_Serialize_Call {
	_c.Call.Return(run)
	return _c
}

// NewResponse creates a new instance of Response. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResponse(t interface {
	mock.TestingT
	Cleanup(func())
}) *Response {
	mock := &Response{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
