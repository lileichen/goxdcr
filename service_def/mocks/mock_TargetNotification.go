// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	base "github.com/couchbase/goxdcr/base"
	mock "github.com/stretchr/testify/mock"
)

// TargetNotification is an autogenerated mock type for the TargetNotification type
type TargetNotification struct {
	mock.Mock
}

type TargetNotification_Expecter struct {
	mock *mock.Mock
}

func (_m *TargetNotification) EXPECT() *TargetNotification_Expecter {
	return &TargetNotification_Expecter{mock: &_m.Mock}
}

// Clone provides a mock function with given fields: numOfReaders
func (_m *TargetNotification) Clone(numOfReaders int) interface{} {
	ret := _m.Called(numOfReaders)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(int) interface{}); ok {
		r0 = rf(numOfReaders)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// TargetNotification_Clone_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Clone'
type TargetNotification_Clone_Call struct {
	*mock.Call
}

// Clone is a helper method to define mock.On call
//   - numOfReaders int
func (_e *TargetNotification_Expecter) Clone(numOfReaders interface{}) *TargetNotification_Clone_Call {
	return &TargetNotification_Clone_Call{Call: _e.mock.On("Clone", numOfReaders)}
}

func (_c *TargetNotification_Clone_Call) Run(run func(numOfReaders int)) *TargetNotification_Clone_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *TargetNotification_Clone_Call) Return(_a0 interface{}) *TargetNotification_Clone_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetNotification_Clone_Call) RunAndReturn(run func(int) interface{}) *TargetNotification_Clone_Call {
	_c.Call.Return(run)
	return _c
}

// GetReplicasInfo provides a mock function with given fields:
func (_m *TargetNotification) GetReplicasInfo() (int, *base.VbHostsMapType, *base.StringStringMap, []uint16) {
	ret := _m.Called()

	var r0 int
	var r1 *base.VbHostsMapType
	var r2 *base.StringStringMap
	var r3 []uint16
	if rf, ok := ret.Get(0).(func() (int, *base.VbHostsMapType, *base.StringStringMap, []uint16)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() *base.VbHostsMapType); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*base.VbHostsMapType)
		}
	}

	if rf, ok := ret.Get(2).(func() *base.StringStringMap); ok {
		r2 = rf()
	} else {
		if ret.Get(2) != nil {
			r2 = ret.Get(2).(*base.StringStringMap)
		}
	}

	if rf, ok := ret.Get(3).(func() []uint16); ok {
		r3 = rf()
	} else {
		if ret.Get(3) != nil {
			r3 = ret.Get(3).([]uint16)
		}
	}

	return r0, r1, r2, r3
}

// TargetNotification_GetReplicasInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetReplicasInfo'
type TargetNotification_GetReplicasInfo_Call struct {
	*mock.Call
}

// GetReplicasInfo is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) GetReplicasInfo() *TargetNotification_GetReplicasInfo_Call {
	return &TargetNotification_GetReplicasInfo_Call{Call: _e.mock.On("GetReplicasInfo")}
}

func (_c *TargetNotification_GetReplicasInfo_Call) Run(run func()) *TargetNotification_GetReplicasInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_GetReplicasInfo_Call) Return(_a0 int, _a1 *base.VbHostsMapType, _a2 *base.StringStringMap, _a3 []uint16) *TargetNotification_GetReplicasInfo_Call {
	_c.Call.Return(_a0, _a1, _a2, _a3)
	return _c
}

func (_c *TargetNotification_GetReplicasInfo_Call) RunAndReturn(run func() (int, *base.VbHostsMapType, *base.StringStringMap, []uint16)) *TargetNotification_GetReplicasInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetTargetBucketInfo provides a mock function with given fields:
func (_m *TargetNotification) GetTargetBucketInfo() base.BucketInfoMapType {
	ret := _m.Called()

	var r0 base.BucketInfoMapType
	if rf, ok := ret.Get(0).(func() base.BucketInfoMapType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.BucketInfoMapType)
		}
	}

	return r0
}

// TargetNotification_GetTargetBucketInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTargetBucketInfo'
type TargetNotification_GetTargetBucketInfo_Call struct {
	*mock.Call
}

// GetTargetBucketInfo is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) GetTargetBucketInfo() *TargetNotification_GetTargetBucketInfo_Call {
	return &TargetNotification_GetTargetBucketInfo_Call{Call: _e.mock.On("GetTargetBucketInfo")}
}

func (_c *TargetNotification_GetTargetBucketInfo_Call) Run(run func()) *TargetNotification_GetTargetBucketInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_GetTargetBucketInfo_Call) Return(_a0 base.BucketInfoMapType) *TargetNotification_GetTargetBucketInfo_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetNotification_GetTargetBucketInfo_Call) RunAndReturn(run func() base.BucketInfoMapType) *TargetNotification_GetTargetBucketInfo_Call {
	_c.Call.Return(run)
	return _c
}

// GetTargetBucketUUID provides a mock function with given fields:
func (_m *TargetNotification) GetTargetBucketUUID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TargetNotification_GetTargetBucketUUID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTargetBucketUUID'
type TargetNotification_GetTargetBucketUUID_Call struct {
	*mock.Call
}

// GetTargetBucketUUID is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) GetTargetBucketUUID() *TargetNotification_GetTargetBucketUUID_Call {
	return &TargetNotification_GetTargetBucketUUID_Call{Call: _e.mock.On("GetTargetBucketUUID")}
}

func (_c *TargetNotification_GetTargetBucketUUID_Call) Run(run func()) *TargetNotification_GetTargetBucketUUID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_GetTargetBucketUUID_Call) Return(_a0 string) *TargetNotification_GetTargetBucketUUID_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetNotification_GetTargetBucketUUID_Call) RunAndReturn(run func() string) *TargetNotification_GetTargetBucketUUID_Call {
	_c.Call.Return(run)
	return _c
}

// GetTargetServerVBMap provides a mock function with given fields:
func (_m *TargetNotification) GetTargetServerVBMap() base.KvVBMapType {
	ret := _m.Called()

	var r0 base.KvVBMapType
	if rf, ok := ret.Get(0).(func() base.KvVBMapType); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(base.KvVBMapType)
		}
	}

	return r0
}

// TargetNotification_GetTargetServerVBMap_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTargetServerVBMap'
type TargetNotification_GetTargetServerVBMap_Call struct {
	*mock.Call
}

// GetTargetServerVBMap is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) GetTargetServerVBMap() *TargetNotification_GetTargetServerVBMap_Call {
	return &TargetNotification_GetTargetServerVBMap_Call{Call: _e.mock.On("GetTargetServerVBMap")}
}

func (_c *TargetNotification_GetTargetServerVBMap_Call) Run(run func()) *TargetNotification_GetTargetServerVBMap_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_GetTargetServerVBMap_Call) Return(_a0 base.KvVBMapType) *TargetNotification_GetTargetServerVBMap_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetNotification_GetTargetServerVBMap_Call) RunAndReturn(run func() base.KvVBMapType) *TargetNotification_GetTargetServerVBMap_Call {
	_c.Call.Return(run)
	return _c
}

// GetTargetStorageBackend provides a mock function with given fields:
func (_m *TargetNotification) GetTargetStorageBackend() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// TargetNotification_GetTargetStorageBackend_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTargetStorageBackend'
type TargetNotification_GetTargetStorageBackend_Call struct {
	*mock.Call
}

// GetTargetStorageBackend is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) GetTargetStorageBackend() *TargetNotification_GetTargetStorageBackend_Call {
	return &TargetNotification_GetTargetStorageBackend_Call{Call: _e.mock.On("GetTargetStorageBackend")}
}

func (_c *TargetNotification_GetTargetStorageBackend_Call) Run(run func()) *TargetNotification_GetTargetStorageBackend_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_GetTargetStorageBackend_Call) Return(_a0 string) *TargetNotification_GetTargetStorageBackend_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetNotification_GetTargetStorageBackend_Call) RunAndReturn(run func() string) *TargetNotification_GetTargetStorageBackend_Call {
	_c.Call.Return(run)
	return _c
}

// IsSourceNotification provides a mock function with given fields:
func (_m *TargetNotification) IsSourceNotification() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TargetNotification_IsSourceNotification_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsSourceNotification'
type TargetNotification_IsSourceNotification_Call struct {
	*mock.Call
}

// IsSourceNotification is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) IsSourceNotification() *TargetNotification_IsSourceNotification_Call {
	return &TargetNotification_IsSourceNotification_Call{Call: _e.mock.On("IsSourceNotification")}
}

func (_c *TargetNotification_IsSourceNotification_Call) Run(run func()) *TargetNotification_IsSourceNotification_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_IsSourceNotification_Call) Return(_a0 bool) *TargetNotification_IsSourceNotification_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TargetNotification_IsSourceNotification_Call) RunAndReturn(run func() bool) *TargetNotification_IsSourceNotification_Call {
	_c.Call.Return(run)
	return _c
}

// Recycle provides a mock function with given fields:
func (_m *TargetNotification) Recycle() {
	_m.Called()
}

// TargetNotification_Recycle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Recycle'
type TargetNotification_Recycle_Call struct {
	*mock.Call
}

// Recycle is a helper method to define mock.On call
func (_e *TargetNotification_Expecter) Recycle() *TargetNotification_Recycle_Call {
	return &TargetNotification_Recycle_Call{Call: _e.mock.On("Recycle")}
}

func (_c *TargetNotification_Recycle_Call) Run(run func()) *TargetNotification_Recycle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TargetNotification_Recycle_Call) Return() *TargetNotification_Recycle_Call {
	_c.Call.Return()
	return _c
}

func (_c *TargetNotification_Recycle_Call) RunAndReturn(run func()) *TargetNotification_Recycle_Call {
	_c.Call.Return(run)
	return _c
}

// NewTargetNotification creates a new instance of TargetNotification. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTargetNotification(t interface {
	mock.TestingT
	Cleanup(func())
}) *TargetNotification {
	mock := &TargetNotification{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}